package domain

import (
	"errors"
	"time"
)

type Enrollment struct {
	ID        int64     `json:"id"`
	StudentID string    `json:"student_id"`
	CourseID  int64     `json:"course_id"`
	CourseName string   `json:"course_name,omitempty"` // Diambil dari Course Service (Tingkat 3)
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type EnrollmentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	CourseID  int64  `json:"course_id" binding:"required"`
}

// Sentinel Errors sesuai kontrak Bagian 6
var (
	ErrValidation          = errors.New("permintaan tidak valid")
	ErrCourseNotFound      = errors.New("course tidak ditemukan")
	ErrEnrollmentNotFound  = errors.New("enrollment tidak ditemukan")
	ErrAlreadyEnrolled     = errors.New("mahasiswa sudah terdaftar di course ini")
	ErrNoSeat              = errors.New("kuota course sudah penuh")
	ErrUpstreamUnavailable = errors.New("layanan course tidak tersedia atau timeout")
	ErrInternal            = errors.New("terjadi kesalahan internal")
)