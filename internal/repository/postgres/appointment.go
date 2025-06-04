package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/pkg/log"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
			id, doctor_id, user_id, schedule_id, status, meeting_url
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
			doctor_id, user_id, schedule_id, status, meeting_url
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id;
	`

	args := []any{
		data.DoctorID,
		data.UserID,
		data.ScheduleID,
		data.Status,
		data.MeetingURL,
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
		id, doctor_id, user_id, schedule_id, status, meeting_url
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
		a.meeting_url,
		d.id AS doctor_id,
		d.name AS doctor_name,
		d.specialization,
		d.phone AS doctor_phone,
		d.gender AS doctor_gender,
		d.visit_type AS doctor_visit_type,
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

func (r *AppointmentRepository) UpdateMeetingURL(ctx context.Context, id string, meetingURL string) error {
	// First, check if the appointment exists and get its status
	var appointment struct {
		ID     string
		Status string
	}
	checkQuery := `SELECT CAST(id AS TEXT) as id, status FROM appointments WHERE id = $1;`
	err := r.db.GetContext(ctx, &appointment, checkQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.LoggerFromContext(ctx).Error("appointment not found", zap.String("id", id))
			return store.ErrorNotFound
		}
		log.LoggerFromContext(ctx).Error("database error", zap.Error(err))
		return err
	}

	// If appointment exists but is not active, return not found
	if appointment.Status != "active" {
		log.LoggerFromContext(ctx).Error("appointment is not active",
			zap.String("id", id),
			zap.String("status", appointment.Status))
		return store.ErrorNotFound
	}

	// Update the meeting URL
	query := `UPDATE appointments SET meeting_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND status = 'active';`
	result, err := r.db.ExecContext(ctx, query, meetingURL, id)
	if err != nil {
		log.LoggerFromContext(ctx).Error("failed to update meeting URL", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.LoggerFromContext(ctx).Error("failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		log.LoggerFromContext(ctx).Error("no rows affected", zap.String("id", id))
		return store.ErrorNotFound
	}

	log.LoggerFromContext(ctx).Info("meeting URL updated successfully",
		zap.String("id", id),
		zap.String("meeting_url", meetingURL))
	return nil
}

func (r *AppointmentRepository) Complete(ctx context.Context, id string) (err error) {
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

	updateAppointment := `UPDATE appointments SET status = 'completed', updated_at = CURRENT_TIMESTAMP WHERE id = $1;`
	if _, err = tx.ExecContext(ctx, updateAppointment, id); err != nil {
		return err
	}

	updateSchedule := `UPDATE schedule SET is_available = TRUE WHERE id = $1;`
	if _, err = tx.ExecContext(ctx, updateSchedule, scheduleID); err != nil {
		return err
	}

	return nil
}
