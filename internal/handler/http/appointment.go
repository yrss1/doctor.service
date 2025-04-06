package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/appointment"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	doctorService *doctorService.Service
}

func NewAppointmentHandler(doctorService doctorService.Service) *AppointmentHandler {
	return &AppointmentHandler{
		doctorService: &doctorService,
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
	res, err := h.doctorService.ListAppointment(c)
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

	res, err := h.doctorService.CreateAppointment(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AppointmentHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorService.GetAppointmentByID(c, id)
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

	if err := h.doctorService.CancelAppointmentByID(c, id); err != nil {
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

	res, err := h.doctorService.ListAppointmentsByUserID(c, id)
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
