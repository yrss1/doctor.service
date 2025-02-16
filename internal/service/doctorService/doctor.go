package doctorService

import (
	"context"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
)

func (s *Service) ListDoctor(ctx context.Context) (res []doctor.Response, err error) {
	data, err := s.doctorRepository.List(ctx)
	if err != nil {
		return
	}

	res = doctor.ParseFromEntities(data)

	return
}
