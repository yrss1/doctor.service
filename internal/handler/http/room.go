package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/room"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	doctorservice *doctorservice.Service
}

func NewRoomHandler(doctorservice doctorservice.Service) *RoomHandler {
	return &RoomHandler{
		doctorservice: &doctorservice,
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

	res, err := h.doctorservice.CreateRoom(c.Request.Context(), req)
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
