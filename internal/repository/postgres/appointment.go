package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/appointment"

	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/jmoiron/sqlx"
)

type AppointmentRepository struct {
	db *sqlx.DB
}

func NewAppointmentRepository(db *sqlx.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) List(ctx context.Context) (dest []appointment.Entity, err error) {
	query := `
		SELECT 
			id, doctor_id, user_id, schedule_id, status
		FROM appointments;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *AppointmentRepository) Add(ctx context.Context, data appointment.Entity) (id string, err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	insertQuery := `
		INSERT INTO appointments (
			doctor_id, user_id, schedule_id, status
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id;
	`

	args := []any{
		data.DoctorID,
		data.UserID,
		data.ScheduleID,
		data.Status,
	}

	if err = tx.QueryRowContext(ctx, insertQuery, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
		return "", err
	}

	updateQuery := `UPDATE schedule SET is_available = FALSE WHERE id = $1;`

	if _, err = tx.ExecContext(ctx, updateQuery, data.ScheduleID); err != nil {
		return "", err
	}

	return id, nil
}

func (r *AppointmentRepository) Get(ctx context.Context, id string) (dest appointment.Entity, err error) {
	query := `
	SELECT 
		id, doctor_id, user_id, schedule_id, status
	FROM appointments
	WHERE id = $1;
	`
	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *AppointmentRepository) Cancel(ctx context.Context, id string) (err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var scheduleID string
	getScheduleQuery := `SELECT schedule_id FROM appointments WHERE id = $1 AND status = 'active';`
	if err = tx.GetContext(ctx, &scheduleID, getScheduleQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrorNotFound
		}
		return err
	}

	updateAppointment := `UPDATE appointments SET status = 'canceled' WHERE id = $1;`
	if _, err = tx.ExecContext(ctx, updateAppointment, id); err != nil {
		return err
	}

	updateSchedule := `UPDATE schedule SET is_available = TRUE WHERE id = $1;`
	if _, err = tx.ExecContext(ctx, updateSchedule, scheduleID); err != nil {
		return err
	}

	return nil
}

func (r *AppointmentRepository) ListByUserID(ctx context.Context, userID string) ([]appointment.EntityView, error) {
	query := `
	SELECT
		a.id AS appointment_id,
		a.status,
		a.created_at,
		a.updated_at,
		d.id AS doctor_id,
		d.name AS doctor_name,
		d.specialization,
		d.phone AS doctor_phone,  -- ← добавлена запятая
		s.slot_start,
		s.slot_end
	FROM appointments a
	JOIN doctors d ON d.id = a.doctor_id
	JOIN schedule s ON s.id = a.schedule_id
	WHERE a.user_id = $1
	ORDER BY a.created_at DESC;
	`

	var result []appointment.EntityView
	if err := r.db.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, err
	}

	return result, nil
}
