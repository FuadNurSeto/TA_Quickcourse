package service

import (
	"context"
	"log"

	"enrollment-service/internal/domain"
)

// Dideklarasikan di sisi Service (Dependency Inversion)[cite: 1]
type EnrollmentRepository interface {
	Create(ctx context.Context, e *domain.Enrollment) error
	GetByID(ctx context.Context, id int64) (domain.Enrollment, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type CourseClient interface {
	GetCourse(ctx context.Context, courseID int64) (string, error)
	Reserve(ctx context.Context, courseID int64) error
	Release(ctx context.Context, courseID int64) error
}

type EnrollmentService struct {
	repo   EnrollmentRepository
	client CourseClient
}

func NewEnrollmentService(repo EnrollmentRepository, client CourseClient) *EnrollmentService {
	return &EnrollmentService{repo: repo, client: client}
}

func (s *EnrollmentService) Enroll(ctx context.Context, req domain.EnrollmentRequest) (domain.Enrollment, error) {
	// 1. Pesan kursi (reserve) ke Course Service via HTTP (BR-01, BR-04, BR-05)[cite: 1]
	err := s.client.Reserve(ctx, req.CourseID)
	if err != nil {
		return domain.Enrollment{}, err
	}

	// 2. Jika reserve sukses, simpan data ke database Enrollment
	e := domain.Enrollment{
		StudentID: req.StudentID,
		CourseID:  req.CourseID,
		Status:    "ACTIVE",
	}
	
	err = s.repo.Create(ctx, &e)
	if err != nil {
		// 3. BR-06: Compensating Transaction[cite: 1]
		// Jika simpan database gagal (misal krn mendaftar 2 kali / BR-02), 
		// kita WAJIB mengembalikan kursi yang terlanjur di-reserve tadi agar kuota tidak hilang[cite: 1].
		log.Printf("Gagal simpan enrollment, memicu kompensasi pelepasan kursi untuk course %d: %v", req.CourseID, err)
		
		// Dijalankan secara asinkron dengan context baru agar tidak terganggu timeout request saat ini
		go func(cID int64) {
			bgCtx := context.Background()
			if relErr := s.client.Release(bgCtx, cID); relErr != nil {
				// Tercatat di log dengan penanda khusus sesuai instruksi BR-06[cite: 1]
				log.Printf("[CRITICAL] Kompensasi release gagal untuk course %d: %v", cID, relErr)
			}
		}(req.CourseID)

		return domain.Enrollment{}, err
	}

	return e, nil
}

func (s *EnrollmentService) Cancel(ctx context.Context, id int64) error {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// BR-03: Mencegah refund kuota ganda (Idempotensi)[cite: 1]
	// Jika status sudah CANCELLED, stop di sini, jangan release kuota lagi[cite: 1].
	if e.Status == "CANCELLED" {
		return nil 
	}

	err = s.repo.UpdateStatus(ctx, id, "CANCELLED")
	if err != nil {
		return err
	}

	// Kembalikan kursi ke Course Service via HTTP[cite: 1]
	return s.client.Release(ctx, e.CourseID)
}

func (s *EnrollmentService) GetByID(ctx context.Context, id int64) (domain.Enrollment, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return e, err
	}

	// Tingkat 3: Mengambil nama course dengan menembak Course Service[cite: 1]
	courseName, err := s.client.GetCourse(ctx, e.CourseID)
	if err == nil {
		e.CourseName = courseName
	} else {
		// Kalau gagal, jangan batalkan response, cukup kasih log saja
		log.Printf("Peringatan: Gagal mengambil nama course: %v", err)
	}

	return e, nil
}