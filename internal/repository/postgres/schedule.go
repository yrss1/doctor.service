package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/yrss1/doctor.service/internal/domain/schedule"

	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) List(ctx context.Context) (dest []schedule.Entity, err error) {
	query := `
		SELECT 
			id, doctor_id, slot_start, slot_end, is_available
		FROM schedule;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *ScheduleRepository) Add(ctx context.Context, data schedule.Entity) (id string, err error) {
	query :=
		`
	INSERT INTO schedule (
		doctor_id, slot_start, slot_end, is_available
	) VALUES (
		$1, $2, $3, $4
	) RETURNING id;

	`

	args := []any{
		data.DoctorID,
		data.SlotStart,
		data.SlotEnd,
		data.IsAvailable,
	}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ScheduleRepository) Get(ctx context.Context, id string) (dest schedule.Entity, err error) {
	query := `
	SELECT 
		id, doctor_id, slot_start, slot_end, is_available
	FROM schedule
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

func (r *ScheduleRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM schedule 
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

func (r *ScheduleRepository) ListByDoctorID(ctx context.Context, doctorID string) (dest []schedule.Entity, err error) {
	query := `
		SELECT 
			id, doctor_id, slot_start, slot_end, is_available
		FROM schedule
		WHERE doctor_id = $1;
		`

	args := []any{doctorID}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}
