package http

import (
	"errors"
	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/internal/service/doctorService"
	"github.com/yrss1/doctor.service/pkg/server/response"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	doctorService *doctorService.Service
}

func NewReviewHandler(doctorService doctorService.Service) *ReviewHandler {
	return &ReviewHandler{
		doctorService: &doctorService,
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
	res, err := h.doctorService.ListReview(c)
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

	res, err := h.doctorService.CreateReview(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ReviewHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.doctorService.GetReviewByID(c, id)
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

	if err := h.doctorService.DeleteReviewByID(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}
}
