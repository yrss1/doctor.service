package http

import (
	"github.com/yrss1/doctor.service/internal/domain/doctor"
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
		api.POST("/", h.Add)
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

func (h *DoctorHandler) Add(c *gin.Context) {
	req := doctor.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorService.CreateDoctor(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
