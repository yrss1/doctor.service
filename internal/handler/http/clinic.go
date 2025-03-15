package http

import (
	"errors"
	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	doctorService *doctorService.Service
}

func NewClinicHandler(doctorService doctorService.Service) *ClinicHandler {
	return &ClinicHandler{
		doctorService: &doctorService,
	}
}

func (h *ClinicHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/clinics")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)
		api.GET("/:id", h.get)
		api.DELETE("/:id", h.delete)
	}
}

func (h *ClinicHandler) list(c *gin.Context) {
	res, err := h.doctorService.ListClinic(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ClinicHandler) add(c *gin.Context) {
	req := clinic.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorService.CreateClinic(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ClinicHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorService.GetClinicByID(c, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

func (h *ClinicHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.doctorService.DeleteClinicByID(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}
