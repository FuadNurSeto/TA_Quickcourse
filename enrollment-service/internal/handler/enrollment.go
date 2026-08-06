package handler

import (
	"net/http"
	"strconv"

	"enrollment-service/internal/domain"
	"enrollment-service/internal/service"
	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	svc *service.EnrollmentService
}

func NewEnrollmentHandler(svc *service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	var req domain.EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, res := MapError(domain.ErrValidation)
		c.JSON(code, res)
		return
	}

	enrollment, err := h.svc.Enroll(c.Request.Context(), req)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Success: true, Data: enrollment})
}

func (h *EnrollmentHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		code, res := MapError(domain.ErrValidation)
		c.JSON(code, res)
		return
	}

	enrollment, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: enrollment})
}

func (h *EnrollmentHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		code, res := MapError(domain.ErrValidation)
		c.JSON(code, res)
		return
	}

	err = h.svc.Cancel(c.Request.Context(), id)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: "pendaftaran berhasil dibatalkan"})
}