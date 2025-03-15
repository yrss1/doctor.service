package doctorService

import (
	"context"
	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListDoctor(ctx context.Context) (res []doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListUsers")

	data, err := s.doctorRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntities(data)

	return
}

func (s *Service) CreateDoctor(ctx context.Context, req doctor.Request) (res doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateUser")
	data := doctor.Entity{
		Name:           req.Name,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Price:          req.Price,
		Rating:         req.Rating,
		Address:        req.Address,
		Phone:          req.Phone,
		ClinicID:       req.ClinicID,
	}

	data.ID, err = s.doctorRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntity(data)

	return
}

func (s *Service) GetDoctorByID(ctx context.Context, id string) (res doctor.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetDoctorByID")

	data, err := s.doctorRepository.Get(ctx, id)
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
