package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListAppointment(ctx context.Context) (res []appointment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListAppointment")

	data, err := s.appointmentRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = appointment.ParseFromEntities(data)

	return
}

func (s *Service) CreateAppointment(ctx context.Context, req appointment.Request) (res appointment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateAppointment")

	data := appointment.Entity{
		DoctorID:   req.DoctorID,
		UserID:     req.UserID,
		ScheduleID: req.ScheduleID,
		Status:     req.Status,
	}

	data.ID, err = s.appointmentRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = appointment.ParseFromEntity(data)

	return
}

func (s *Service) GetAppointmentByID(ctx context.Context, id string) (res appointment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetAppointmentByID")

	data, err := s.appointmentRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = appointment.ParseFromEntity(data)

	return
}

func (s *Service) CancelAppointmentByID(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("CancelAppointmentByID")

	err = s.appointmentRepository.Cancel(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
	}

	return
}

func (s *Service) ListAppointmentsByUserID(ctx context.Context, id string) (data []appointment.EntityView, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListAppointmentByUserID")

	data, err = s.appointmentRepository.ListByUserID(ctx, id)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	return
}

func (s *Service) UpdateAppointmentMeetingURL(ctx context.Context, id string, meetingURL string) error {
	logger := log.LoggerFromContext(ctx).Named("UpdateAppointmentMeetingURL")

	if err := s.appointmentRepository.UpdateMeetingURL(ctx, id, meetingURL); err != nil {
		logger.Error("failed to update meeting URL", zap.Error(err))
		return err
	}

	return nil
}
