package handler

import "course-service/internal/domain"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// Fungsi helper untuk mapping error ke kode HTTP dan pesan JSON
func MapError(err error) (int, ErrorResponse) {
	var httpCode int
	var errCode string
	var message = err.Error()

	switch err {
	case domain.ErrCourseNotFound:
		httpCode = 404
		errCode = "ERR_COURSE_NOT_FOUND"
	case domain.ErrDuplicateCode:
		httpCode = 409
		errCode = "ERR_ALREADY_EXISTS" // Custom
	case domain.ErrNoSeat:
		httpCode = 409
		errCode = "ERR_NO_SEAT"
	case domain.ErrInvalidRequest:
		httpCode = 400
		errCode = "ERR_VALIDATION"
	default:
		httpCode = 500
		errCode = "ERR_INTERNAL"
		message = "terjadi kesalahan pada server" // Sembunyikan detail dari client (Keamanan)
	}

	return httpCode, ErrorResponse{
		Success: false,
		Error:   ErrorDetail{Code: errCode, Message: message},
	}
}