package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListDoctorWithSchedules(ctx context.Context) (res []doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListUsers")

	data, err := s.doctorRepository.ListWithSchedules(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntities(data)

	return
}

func (s *Service) GetDoctorByIDWithSchedules(ctx context.Context, id string) (res doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetDoctorByID")

	data, err := s.doctorRepository.GetWithSchedules(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntity(data)

	return
}

func (s *Service) DeleteDoctorByID(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteDoctorByID")

	err = s.doctorRepository.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
	}

	return
}

func (s *Service) SearchWithSchedules(ctx context.Context, req doctor.Request) (res []doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchWithSchedules")

	if req.Name != nil {
		logger = logger.With(zap.String("name", *(req.Name)))
	}
	if req.Specialization != nil {
		logger = logger.With(zap.String("specialization", *(req.Specialization)))
	}
	if req.ClinicName != nil {
		logger = logger.With(zap.String("clinicName", *(req.ClinicName)))
	}
	searchData := doctor.Entity{
		Name:           req.Name,
		Specialization: req.Specialization,
		ClinicName:     req.ClinicName,
	}
	data, err := s.doctorRepository.SearchWithSchedules(ctx, searchData)
	if err != nil {
		logger.Error("failed to search doctors", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntities(data)

	return
}
