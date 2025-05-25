package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/internal/provider/meet"
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

	// Create appointment entity
	data := appointment.Entity{
		DoctorID:   req.DoctorID,
		UserID:     req.UserID,
		ScheduleID: req.ScheduleID,
		Status:     req.Status,
	}

	// Add appointment to database
	data.ID, err = s.appointmentRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create appointment", zap.Error(err))
		return
	}

	// Get schedule information for meeting times
	schedule, err := s.scheduleRepository.Get(ctx, *req.ScheduleID)
	if err != nil {
		logger.Error("failed to get schedule", zap.Error(err))
		return
	}

	// Get doctor information for email
	doctor, err := s.doctorRepository.GetWithSchedules(ctx, *req.DoctorID)
	if err != nil {
		logger.Error("failed to get doctor", zap.Error(err))
		return
	}

	// Create meeting request
	meetReq := meet.Request{
		UserEmail:   *req.UserID + "@gmail.com", // Note: You'll need to modify this based on how you store user emails
		DoctorEmail: *doctor.Phone,              // Using phone as email since that's what's stored
		StartTime:   schedule.SlotStart.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:     schedule.SlotEnd.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Create Google Meet meeting
	_, err = s.CreateMeeting(ctx, meetReq)
	if err != nil {
		logger.Error("failed to create meeting", zap.Error(err))
		// Note: We don't return here as the appointment is already created
		// You might want to add a field to the appointment to track meeting creation status
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
