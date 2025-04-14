package http

import (
	"fmt"
	"net/http"

	"github.com/yrss1/doctor.service/internal/provider/meet"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"

	"github.com/gin-gonic/gin"
)

type MeetHandler struct {
	doctorservice *doctorservice.Service
}

func NewMeetHandler(doctorservice doctorservice.Service) *MeetHandler {
	return &MeetHandler{
		doctorservice: &doctorservice,
	}
}

func (h *MeetHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/meets")
	{
		api.GET("/login", h.login)
		api.POST("/create", h.createMeeting)
		api.GET("/oauth2callback", h.oauthCallback)

	}
}

func (h *MeetHandler) login(c *gin.Context) {
	link, err := h.doctorservice.Login(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	c.Redirect(http.StatusFound, link)
}

func (h *MeetHandler) createMeeting(c *gin.Context) {
	var req meet.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	res, err := h.doctorservice.CreateMeeting(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *MeetHandler) oauthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, fmt.Errorf("missing code in query"), code)
		return
	}
	token, err := h.doctorservice.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, token)
}
