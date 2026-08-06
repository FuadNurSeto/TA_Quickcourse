package handler

import (
	"net/http"
	"strconv"
	"course-service/internal/domain"
	"course-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	svc *service.CourseService
}

func NewCourseHandler(svc *service.CourseService) *CourseHandler {
	return &CourseHandler{svc: svc}
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req domain.CourseRequest
	// Validasi body JSON[cite: 1]
	if err := c.ShouldBindJSON(&req); err != nil {
		code, res := MapError(domain.ErrInvalidRequest)
		c.JSON(code, res)
		return
	}

	course, err := h.svc.CreateCourse(c.Request.Context(), req)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Success: true, Data: course})
}

func (h *CourseHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		code, res := MapError(domain.ErrInvalidRequest)
		c.JSON(code, res)
		return
	}

	course, err := h.svc.GetCourse(c.Request.Context(), id)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: course})
}

// Endpoint internal untuk Enrollment Service (BR-01)[cite: 1]
func (h *CourseHandler) Reserve(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Reserve(c.Request.Context(), id); err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: "reserved"})
}

// Endpoint internal kompensasi pembatalan (BR-03)[cite: 1]
func (h *CourseHandler) Release(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Release(c.Request.Context(), id); err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: "released"})
}
// Endpoint untuk cek sisa kursi[cite: 1]
func (h *CourseHandler) GetAvailability(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		code, res := MapError(domain.ErrInvalidRequest)
		c.JSON(code, res)
		return
	}

	course, err := h.svc.GetCourse(c.Request.Context(), id)
	if err != nil {
		code, res := MapError(err)
		c.JSON(code, res)
		return
	}

	// Mengembalikan capacity, taken, dan remaining sesuai kontrak Bagian 4.1[cite: 1]
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: gin.H{
			"capacity":  course.Capacity,
			"taken":     course.Taken,
			"remaining": course.Remaining,
		},
	})
}