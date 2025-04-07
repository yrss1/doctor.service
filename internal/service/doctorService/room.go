package doctorService

import (
	"context"
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/room"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) CreateRoom(ctx context.Context, req room.Entity) (id string, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateRoom")

	appt, err := s.GetAppointmentByID(ctx, req.AppointmentID)
	if err != nil {
		logger.Error("appointment not found", zap.Error(err))
		return id, errors.New("access denied")
	}

	if req.DoctorID != *appt.DoctorID || req.UserID != *appt.UserID {
		logger.Error("acces denied", zap.Error(err))
		return id, errors.New("access denied")
	}
	id += "room-" + appt.ID

	return
}
