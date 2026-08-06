package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"enrollment-service/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type enrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *enrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, e *domain.Enrollment) error {
	query := `INSERT INTO enrollments (student_id, course_id, status) 
	          VALUES ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, e.StudentID, e.CourseID, e.Status).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		// BR-02: Menangkap error constraint unik dari PostgreSQL
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyEnrolled
		}
		return fmt.Errorf("gagal insert enrollment: %w", err)
	}
	return nil
}

func (r *enrollmentRepository) GetByID(ctx context.Context, id int64) (domain.Enrollment, error) {
	query := `SELECT id, student_id, course_id, status, created_at FROM enrollments WHERE id = $1`
	var e domain.Enrollment
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.Status, &e.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, domain.ErrEnrollmentNotFound
		}
		return e, fmt.Errorf("gagal get enrollment: %w", err)
	}
	return e, nil
}

func (r *enrollmentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE enrollments SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("gagal update status: %w", err)
	}
	return nil
}