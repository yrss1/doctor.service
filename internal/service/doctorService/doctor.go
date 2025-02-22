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
		Name:       req.Name,
		Specialty:  req.Specialty,
		Experience: req.Experience,
		Price:      req.Price,
		Address:    req.Address,
		ClinicName: req.ClinicName,
		Phone:      req.Phone,
		Email:      req.Email,
		PhotoURL:   req.PhotoURL,
		Education:  req.Education,
	}

	data.ID, err = s.doctorRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = doctor.ParseFromEntity(data)

	return
}
