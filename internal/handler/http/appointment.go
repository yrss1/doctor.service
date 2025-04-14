package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	doctorservice *doctorservice.Service
}

func NewAppointmentHandler(doctorservice doctorservice.Service) *AppointmentHandler {
	return &AppointmentHandler{
		doctorservice: &doctorservice,
	}
}

func (h *AppointmentHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/appointments")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)
		api.GET("/:id", h.get)
		api.GET("/cancel/:id", h.cancel)
		api.GET("/user/:id", h.listByUserID)
	}
}

func (h *AppointmentHandler) list(c *gin.Context) {
	res, err := h.doctorservice.ListAppointment(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AppointmentHandler) add(c *gin.Context) {
	req := appointment.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorservice.CreateAppointment(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AppointmentHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.GetAppointmentByID(c.Request.Context(), id)
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

func (h *AppointmentHandler) cancel(c *gin.Context) {
	id := c.Param("id")

	if err := h.doctorservice.CancelAppointmentByID(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}

func (h *AppointmentHandler) listByUserID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.ListAppointmentsByUserID(c.Request.Context(), id)
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
