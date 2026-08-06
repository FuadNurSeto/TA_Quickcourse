package service

import (
	"context"
	"course-service/internal/domain"
)

// Dideklarasikan di sisi yang memakainya[cite: 1]
type CourseRepository interface {
	Create(ctx context.Context, c *domain.Course) error
	GetByID(ctx context.Context, id int64) (domain.Course, error)
	ReserveSeat(ctx context.Context, id int64) error
	ReleaseSeat(ctx context.Context, id int64) error
}

type CourseService struct {
	repo CourseRepository
}

// Dependency Injection[cite: 1]
func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateCourse(ctx context.Context, req domain.CourseRequest) (domain.Course, error) {
	c := domain.Course{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Lecturer:    req.Lecturer,
		Capacity:    req.Capacity,
	}
	err := s.repo.Create(ctx, &c)
	return c, err
}

func (s *CourseService) GetCourse(ctx context.Context, id int64) (domain.Course, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CourseService) Reserve(ctx context.Context, id int64) error {
	return s.repo.ReserveSeat(ctx, id)
}

func (s *CourseService) Release(ctx context.Context, id int64) error {
	return s.repo.ReleaseSeat(ctx, id)
}