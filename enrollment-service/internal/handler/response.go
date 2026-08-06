package handler

import "enrollment-service/internal/domain"

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

// Pemetaan error mutlak sesuai Bagian 6 di PDF
func MapError(err error) (int, ErrorResponse) {
	var httpCode int
	var errCode string
	var message = err.Error()

	switch err {
	case domain.ErrValidation:
		httpCode = 400
		errCode = "ERR_VALIDATION"
	case domain.ErrCourseNotFound:
		httpCode = 404
		errCode = "ERR_COURSE_NOT_FOUND"
	case domain.ErrEnrollmentNotFound:
		httpCode = 404
		errCode = "ERR_ENROLLMENT_NOT_FOUND"
	case domain.ErrAlreadyEnrolled:
		httpCode = 409
		errCode = "ERR_ALREADY_ENROLLED"
	case domain.ErrNoSeat:
		httpCode = 409
		errCode = "ERR_NO_SEAT"
	case domain.ErrUpstreamUnavailable:
		httpCode = 503
		errCode = "ERR_UPSTREAM_UNAVAILABLE"
	default:
		httpCode = 500
		errCode = "ERR_INTERNAL"
		message = "terjadi kesalahan pada server"
	}

	return httpCode, ErrorResponse{
		Success: false,
		Error:   ErrorDetail{Code: errCode, Message: message},
	}
}