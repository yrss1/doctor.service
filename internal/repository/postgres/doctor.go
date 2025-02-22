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
			id, name, specialty, experience, price, address, clinic_name, 
			phone, email, photo_url, education, rating, reviews_count, is_active
		FROM doctors;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *DoctorRepository) Add(ctx context.Context, data doctor.Entity) (id string, err error) {
	query :=
		`
	INSERT INTO doctors (
		name, specialty, experience, price, address, clinic_name, 
		phone, email, photo_url, education
	) VALUES (
		$1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10
	) RETURNING id;

	`

	args := []any{
		data.Name,
		data.Specialty,
		data.Experience,
		data.Price,
		data.Address,
		data.ClinicName,
		data.Phone,
		data.Email,
		data.PhotoURL,
		data.Education,
	}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}
