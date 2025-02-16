package app

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := servers.Run(); err != nil {
			logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		}
	}()

	logger.Info("HTTP server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", 15*time.Second, "The duration for which the server will wait for existing connections to finish before shutting down")
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("Shutdown initiated...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := servers.Stop(ctx); err != nil {
			logger.Error("ERR_SHUTDOWN_SERVERS", zap.Error(err))
		}
	}()

	wg.Wait()

	logger.Info("All services stopped. Exiting.")
}
