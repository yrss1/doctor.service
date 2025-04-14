package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListClinic(ctx context.Context) (res []clinic.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListClinic")

	data, err := s.clinicRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = clinic.ParseFromEntities(data)

	return
}

func (s *Service) CreateClinic(ctx context.Context, req clinic.Request) (res clinic.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateClinic")

	data := clinic.Entity{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	data.ID, err = s.clinicRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = clinic.ParseFromEntity(data)

	return
}

func (s *Service) GetClinicByID(ctx context.Context, id string) (res clinic.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetClinicByID")

	data, err := s.clinicRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = clinic.ParseFromEntity(data)

	return
}

func (s *Service) DeleteClinicByID(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteClinicByID")

	err = s.clinicRepository.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
	}

	return
}
