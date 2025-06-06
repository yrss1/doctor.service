package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorservice *doctorservice.Service
}

func NewDoctorHandler(doctorservice doctorservice.Service) *DoctorHandler {
	return &DoctorHandler{
		doctorservice: &doctorservice,
	}
}

func (h *DoctorHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/doctors")
	{
		api.GET("/", h.listWithSchedules)
		api.GET("/:id", h.get)
		api.DELETE("/:id", h.delete)
		api.GET("/search", h.search)
	}
}

func (h *DoctorHandler) listWithSchedules(c *gin.Context) {
	res, err := h.doctorservice.ListDoctorWithSchedules(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *DoctorHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.GetDoctorByIDWithSchedules(c.Request.Context(), id)
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

func (h *DoctorHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.doctorservice.DeleteDoctorByID(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}

func (h *DoctorHandler) search(c *gin.Context) {
	req := doctor.Request{
		Name:           ptr(c.Query("name")),
		Specialization: ptr(c.Query("specialization")),
		ClinicName:     ptr(c.Query("clinic_name")),
	}

	res, err := h.doctorservice.SearchWithSchedules(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func ptr[T any](v T) *T {
	return &v
}
