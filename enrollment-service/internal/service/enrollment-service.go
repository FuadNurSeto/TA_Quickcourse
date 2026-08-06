package service

import (
	"context"
	"testing"
	"errors"

	"enrollment-service/internal/domain"
)

// --- MOCK REPOSITORY ---
type mockRepo struct {
	createErr error
}
func (m *mockRepo) Create(ctx context.Context, e *domain.Enrollment) error { return m.createErr }
func (m *mockRepo) GetByID(ctx context.Context, id int64) (domain.Enrollment, error) { return domain.Enrollment{}, nil }
func (m *mockRepo) UpdateStatus(ctx context.Context, id int64, status string) error { return nil }

// --- MOCK HTTP CLIENT ---
type mockClient struct {
	reserveErr error
}
func (m *mockClient) GetCourse(ctx context.Context, id int64) (string, error) { return "Go Dasar", nil }
func (m *mockClient) Reserve(ctx context.Context, id int64) error { return m.reserveErr }
func (m *mockClient) Release(ctx context.Context, id int64) error { return nil }

// --- TABLE-DRIVEN TEST ---
func TestEnrollmentService_Enroll(t *testing.T) {
	// Skenario pengujian wajib dari dokumen
	tests := []struct {
		name        string
		reserveErr  error
		createErr   error
		expectedErr error
	}{
		{
			name:        "Kasus Sukses: Berhasil mendaftar",
			reserveErr:  nil,
			createErr:   nil,
			expectedErr: nil,
		},
		{
			name:        "Kasus Kuota Habis: Course Service menolak reserve",
			reserveErr:  domain.ErrNoSeat,
			createErr:   nil,
			expectedErr: domain.ErrNoSeat,
		},
		{
			name:        "Kasus Course Service Gagal Merespons / Timeout",
			reserveErr:  domain.ErrUpstreamUnavailable,
			createErr:   nil,
			expectedErr: domain.ErrUpstreamUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{createErr: tt.createErr}
			client := &mockClient{reserveErr: tt.reserveErr}
			svc := NewEnrollmentService(repo, client)

			req := domain.EnrollmentRequest{StudentID: "123", CourseID: 1}
			_, err := svc.Enroll(context.Background(), req)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}