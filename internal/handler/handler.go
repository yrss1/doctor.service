package handler

import (
	"github.com/yrss1/doctor.service/internal/config"
	"github.com/yrss1/doctor.service/internal/handler/http"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/server/router"

	"github.com/gin-contrib/timeout"
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
		h.HTTP.Use(timeout.New(
			timeout.WithTimeout(h.dependencies.Configs.APP.Timeout),
			timeout.WithHandler(func(ctx *gin.Context) {
				ctx.Next()
			}),
			timeout.WithResponse(func(ctx *gin.Context) {
				response.StatusRequestTimeout(ctx)
			}),
		))

		doctorHandler := http.NewDoctorHandler(h.dependencies.DoctorService)
		clinicHandler := http.NewClinicHandler(h.dependencies.DoctorService)

		api := h.HTTP.Group("/api/v1")
		{
			doctorHandler.Routes(api)
			clinicHandler.Routes(api)
		}
		return
	}
}
