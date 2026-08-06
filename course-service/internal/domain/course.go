package domain

import "errors"

// Model Database
type Course struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Lecturer    string `json:"lecturer"`
	Capacity    int    `json:"capacity"`
	Taken       int    `json:"taken"`
	Remaining   int    `json:"remaining"` // Dihitung manual: Capacity - Taken
}

// DTO untuk validasi input JSON via Gin
type CourseRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Lecturer    string `json:"lecturer" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,gt=0"`
}

// Sentinel Errors sesuai kontrak dokumen Bagian 6
var (
	ErrCourseNotFound = errors.New("course tidak ditemukan")
	ErrDuplicateCode  = errors.New("kode course sudah dipakai")
	ErrNoSeat         = errors.New("kuota course sudah penuh")
	ErrInvalidRequest = errors.New("permintaan tidak valid")
)