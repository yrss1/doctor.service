package doctorservice

import (
	"context"

	"github.com/yrss1/doctor.service/internal/provider/meet"
	"golang.org/x/oauth2"
)

func (s *Service) Login(ctx context.Context) (string, error) {
	return s.meetClient.LoginURL(), nil
}

func (s *Service) CreateMeeting(ctx context.Context, req meet.Request) (res meet.Response, err error) {
	return s.meetClient.CreateMeetEvent(ctx, req)
}

func (s *Service) ExchangeCode(ctx context.Context, code string) (token *oauth2.Token, err error) {
	return s.meetClient.ExchangeCode(ctx, code)
}
