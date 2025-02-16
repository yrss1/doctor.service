package http

import (
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorService *doctorService.Service
}

func NewDoctorHandler(doctorService doctorService.Service) *DoctorHandler {
	return &DoctorHandler{
		doctorService: &doctorService,
	}
}

func (h *DoctorHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/doctors")
	{
		api.GET("/", h.list)

	}
}

func (h *DoctorHandler) list(c *gin.Context) {
	res, err := h.doctorService.ListDoctor(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
