package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/schedule"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListSchedule(ctx context.Context) (res []schedule.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListSchedules")

	data, err := s.scheduleRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = schedule.ParseFromEntities(data)

	return
}

func (s *Service) CreateSchedule(ctx context.Context, req schedule.Request) (res schedule.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateUser")
	data := schedule.Entity{
		DoctorID:    req.DoctorID,
		SlotStart:   req.SlotStart,
		SlotEnd:     req.SlotEnd,
		IsAvailable: req.IsAvailable,
	}

	data.ID, err = s.scheduleRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = schedule.ParseFromEntity(data)

	return
}

func (s *Service) GetScheduleByID(ctx context.Context, id string) (res schedule.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetScheduleByID")

	data, err := s.scheduleRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = schedule.ParseFromEntity(data)

	return
}

func (s *Service) DeleteScheduleByID(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteScheduleByID")

	err = s.scheduleRepository.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
	}

	return
}

func (s *Service) ListScheduleByDoctorID(ctx context.Context, doctorID string) (res []schedule.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListScheduleByDoctorID")

	data, err := s.scheduleRepository.ListByDoctorID(ctx, doctorID)
	if err != nil {
		logger.Error("failed to select by doctor_id", zap.Error(err))
		return
	}

	res = schedule.ParseFromEntities(data)

	return
}
