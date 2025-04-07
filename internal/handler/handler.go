package handler

import (
	"github.com/yrss1/doctor.service/internal/config"
	"github.com/yrss1/doctor.service/internal/handler/http"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/router"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Configs config.Configs

	DoctorService doctorService.Service
}

type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
}

type Configuration func(h *Handler) error

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()

		doctorHandler := http.NewDoctorHandler(h.dependencies.DoctorService)
		clinicHandler := http.NewClinicHandler(h.dependencies.DoctorService)
		scheduleHandler := http.NewScheduleHandler(h.dependencies.DoctorService)
		appointmentHandler := http.NewAppointmentHandler(h.dependencies.DoctorService)
		reviewHandler := http.NewReviewHandler(h.dependencies.DoctorService)
		roomHanlder := http.NewRoomHandler(h.dependencies.DoctorService)

		api := h.HTTP.Group("/api/v1")
		{
			doctorHandler.Routes(api)
			clinicHandler.Routes(api)
			scheduleHandler.Routes(api)
			appointmentHandler.Routes(api)
			reviewHandler.Routes(api)
			roomHanlder.Routes(api)
		}
		api.GET("/health")

		return
	}
}
