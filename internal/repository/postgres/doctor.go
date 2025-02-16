package postgres

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/doctor"

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
		SELECT id, name, specialty, experience, price 
		FROM doctors
		ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}
