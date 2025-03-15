package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/pkg/store"
)

type ClinicRepository struct {
	db *sqlx.DB
}

func NewClinicRepository(db *sqlx.DB) *ClinicRepository {
	return &ClinicRepository{db: db}
}

func (r *ClinicRepository) List(ctx context.Context) (dest []clinic.Entity, err error) {
	query := `
		SELECT 
			id, name, address, phone
		FROM clinic;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *ClinicRepository) Add(ctx context.Context, data clinic.Entity) (id string, err error) {
	query :=
		`
	INSERT INTO clinic (
		name, address, phone
	) VALUES (
		$1, $2, $3
	) RETURNING id;

	`

	args := []any{
		data.Name,
		data.Address,
		data.Phone,
	}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ClinicRepository) Get(ctx context.Context, id string) (dest clinic.Entity, err error) {
	query := `
	SELECT 
		id, name, address, phone
	FROM clinic
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

func (r *ClinicRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM clinic 
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
