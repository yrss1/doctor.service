package http

import (
	"errors"

	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/internal/service/doctorservice"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	doctorservice *doctorservice.Service
}

func NewReviewHandler(doctorservice doctorservice.Service) *ReviewHandler {
	return &ReviewHandler{
		doctorservice: &doctorservice,
	}
}

func (h *ReviewHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/reviews")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)
		api.GET("/:id", h.get)
		api.DELETE("/:id", h.delete)
	}
}

func (h *ReviewHandler) list(c *gin.Context) {
	res, err := h.doctorservice.ListReview(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ReviewHandler) add(c *gin.Context) {
	req := review.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.doctorservice.CreateReview(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ReviewHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorservice.GetReviewByID(c, id)
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

func (h *ReviewHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.doctorservice.DeleteReviewByID(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}
