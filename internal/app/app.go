package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yrss1/doctor.service/internal/config"
	"github.com/yrss1/doctor.service/internal/handler"
	"github.com/yrss1/doctor.service/internal/repository"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/log"
	"github.com/yrss1/doctor.service/pkg/server"
	"go.uber.org/zap"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

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

	doctorService, err := doctorService.New(
		doctorService.WithDoctorRepository(repositories.Doctor),
		doctorService.WithClinicRepository(repositories.Clinic),
	)
	if err != nil {
		logger.Error("ERR_INIT_DOCTOR_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:       *configs,
			DoctorService: *doctorService,
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
