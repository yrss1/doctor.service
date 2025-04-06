package repository

import (
	"fmt"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/internal/domain/schedule"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/internal/repository/postgres"
	"github.com/yrss1/doctor.service/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	Doctor      doctor.Repository
	Clinic      clinic.Repository
	Schedule    schedule.Repository
	Appointment appointment.Repository
	Review      review.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(dbName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.New(dbName)
		if err != nil {
			return fmt.Errorf("failed to initialize database connection: %w", err)
		}
		if err = store.Migrate(dbName); err != nil {
			return fmt.Errorf("failed to run database migrations: %w", err)
		}

		r.Doctor = postgres.NewDoctorRepository(r.postgres.Client)
		r.Clinic = postgres.NewClinicRepository(r.postgres.Client)
		r.Schedule = postgres.NewScheduleRepository(r.postgres.Client)
		r.Appointment = postgres.NewAppointmentRepository(r.postgres.Client)
		r.Review = postgres.NewReviewRepository(r.postgres.Client)

		return
	}
}
