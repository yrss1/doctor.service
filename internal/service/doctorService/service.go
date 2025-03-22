package doctorService

import (
	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/internal/domain/schedule"
)

type Configuration func(s *Service) error

type Service struct {
	doctorRepository      doctor.Repository
	clinicRepository      clinic.Repository
	scheduleRepository    schedule.Repository
	appointmentRepository appointment.Repository
	reviewRepository      review.Repository
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithDoctorRepository(doctorRepository doctor.Repository) Configuration {
	return func(s *Service) error {
		s.doctorRepository = doctorRepository
		return nil
	}
}

func WithClinicRepository(clinicRepository clinic.Repository) Configuration {
	return func(s *Service) error {
		s.clinicRepository = clinicRepository
		return nil
	}
}

func WithScheduleRepository(scheduleRepository schedule.Repository) Configuration {
	return func(s *Service) error {
		s.scheduleRepository = scheduleRepository
		return nil
	}
}

func WithAppointmentRepository(appointmentRepository appointment.Repository) Configuration {
	return func(s *Service) error {
		s.appointmentRepository = appointmentRepository
		return nil
	}
}

func WithReviewRepository(reviewRepository review.Repository) Configuration {
	return func(s *Service) error {
		s.reviewRepository = reviewRepository
		return nil
	}
}
