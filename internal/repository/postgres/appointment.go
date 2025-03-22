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
	query :=
		`
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

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
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

func (r *AppointmentRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM appointments 
	WHERE id = $1;
	RETURNING id;
	`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}
