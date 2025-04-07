package http

import (
	"github.com/yrss1/doctor.service/internal/domain/room"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	doctorService *doctorService.Service
}

func NewRoomHandler(doctorService doctorService.Service) *RoomHandler {
	return &RoomHandler{
		doctorService: &doctorService,
	}
}

func (h *RoomHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/rooms")
	{
		api.POST("/", h.add)
	}
}

func (h *RoomHandler) add(c *gin.Context) {
	req := room.Entity{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorService.CreateRoom(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
