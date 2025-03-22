package doctorService

import (
	"context"
	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/pkg/log"
	"go.uber.org/zap"
)

func (s *Service) ListReview(ctx context.Context) (res []review.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListReview")

	data, err := s.reviewRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = review.ParseFromEntities(data)

	return
}

func (s *Service) CreateReview(ctx context.Context, req review.Request) (res review.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateReview")

	data := review.Entity{
		DoctorID: req.DoctorID,
		UserID:   req.UserID,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}

	data.ID, err = s.reviewRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = review.ParseFromEntity(data)

	return
}

func (s *Service) GetReviewByID(ctx context.Context, id string) (res review.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetReviewByID")

	data, err := s.reviewRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = review.ParseFromEntity(data)

	return
}

func (s *Service) DeleteReviewByID(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteReviewByID")

	err = s.reviewRepository.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
	}

	return
}
