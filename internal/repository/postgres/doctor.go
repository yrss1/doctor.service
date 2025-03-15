package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/jmoiron/sqlx"
)

type DoctorRepository struct {
	db *sqlx.DB
}

func NewDoctorRepository(db *sqlx.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) List(ctx context.Context) (dest []doctor.Entity, err error) {
	query := `
		SELECT 
			id, name, specialization, experience, price, rating, address, phone, clinic_id
		FROM doctors;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *DoctorRepository) Add(ctx context.Context, data doctor.Entity) (id string, err error) {
	query :=
		`
	INSERT INTO doctors (
		name, specialization, experience, price, rating, address, phone, clinic_id
	) VALUES (
		$1, $2, $3, $4, $5, $6, 
		$7, $8
	) RETURNING id;

	`

	args := []any{
		data.Name,
		data.Specialization,
		data.Experience,
		data.Price,
		data.Rating,
		data.Address,
		data.Phone,
		data.ClinicID,
	}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *DoctorRepository) Get(ctx context.Context, id string) (dest doctor.Entity, err error) {
	query := `
	SELECT 
		id, name, specialization, experience, price, rating, address, phone, clinic_id
	FROM doctors
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

func (r *DoctorRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM doctors 
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
