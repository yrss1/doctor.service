package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yrss1/doctor.service/internal/config"
	"github.com/yrss1/doctor.service/internal/handler"
	"github.com/yrss1/doctor.service/internal/provider/meet"
	"github.com/yrss1/doctor.service/internal/repository"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/log"
	"github.com/yrss1/doctor.service/pkg/server"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	b, err := os.ReadFile("/etc/secrets/credentials.json")
	if err != nil {
		logger.Error("ERR_READ_CLIENT_FILE", zap.Error(err))
		return
	}

	oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		logger.Error("ERR_LOAD_OAUTH_CONFIGS", zap.Error(err))
		return
	}

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	repositories, err := repository.New(repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}

	// Create token storage
	tokenStorage := meet.NewFileTokenStorage("tokens/google_calendar_token.json")

	meetClient, err := meet.New(meet.Credentials{
		URL:          configs.APP.Mode,
		OauthConfig:  oauthConfig,
		OauthToken:   nil, // Will be loaded from storage if available
		TokenStorage: tokenStorage,
	})
	if err != nil {
		logger.Error("ERR_INIT_MEET_CLIENT", zap.Error(err))
		return
	}

	doctorservice, err := doctorservice.New(
		doctorservice.WithDoctorRepository(repositories.Doctor),
		doctorservice.WithClinicRepository(repositories.Clinic),
		doctorservice.WithScheduleRepository(repositories.Schedule),
		doctorservice.WithAppointmentRepository(repositories.Appointment),
		doctorservice.WithReviewRepository(repositories.Review),
		doctorservice.WithMeetClient(*meetClient),
	)
	if err != nil {
		logger.Error("ERR_INIT_DOCTOR_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:       *configs,
			DoctorService: *doctorservice,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port),
	)
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	logger.Info("HTTP server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	errChan := make(chan error, 1)
	go func() {
		errChan <- servers.Run()
	}()

	select {
	case <-quit:
		logger.Info("Shutdown initiated...")
	case err := <-errChan:
		if err != nil {
			logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := servers.Stop(ctx); err != nil {
		logger.Error("ERR_SHUTDOWN_SERVERS", zap.Error(err))
	}

	logger.Info("All services stopped. Exiting.")
}
