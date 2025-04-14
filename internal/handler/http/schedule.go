package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/schedule"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	doctorservice *doctorservice.Service
}

func NewScheduleHandler(doctorservice doctorservice.Service) *ScheduleHandler {
	return &ScheduleHandler{
		doctorservice: &doctorservice,
	}
}

func (h *ScheduleHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/schedules")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)
		api.GET("/:id", h.get)
		api.DELETE("/:id", h.delete)
		api.GET("/byDoctorID/:doctorID", h.listByDoctorID)
	}
}

func (h *ScheduleHandler) list(c *gin.Context) {
	res, err := h.doctorservice.ListSchedule(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ScheduleHandler) add(c *gin.Context) {
	req := schedule.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorservice.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ScheduleHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.GetScheduleByID(c.Request.Context(), id)
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

func (h *ScheduleHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.doctorservice.DeleteScheduleByID(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}

func (h *ScheduleHandler) listByDoctorID(c *gin.Context) {
	doctorID := c.Param("doctorID")
	res, err := h.doctorservice.ListScheduleByDoctorID(c.Request.Context(), doctorID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
