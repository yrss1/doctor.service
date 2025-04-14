package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/clinic"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	doctorservice *doctorservice.Service
}

func NewClinicHandler(doctorservice doctorservice.Service) *ClinicHandler {
	return &ClinicHandler{
		doctorservice: &doctorservice,
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
	res, err := h.doctorservice.ListClinic(c.Request.Context())
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

	res, err := h.doctorservice.CreateClinic(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ClinicHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.GetClinicByID(c.Request.Context(), id)
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

	if err := h.doctorservice.DeleteClinicByID(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}
