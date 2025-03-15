package doctorService

import (
	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/internal/domain/doctor"
)

type Configuration func(s *Service) error

type Service struct {
	doctorRepository doctor.Repository
	clinicRepository clinic.Repository
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
