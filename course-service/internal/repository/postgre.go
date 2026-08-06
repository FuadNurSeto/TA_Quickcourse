package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"course-service/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type courseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *courseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, c *domain.Course) error {
	query := `INSERT INTO courses (code, name, description, lecturer, capacity) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, taken`
	
	err := r.db.QueryRowContext(ctx, query, c.Code, c.Name, c.Description, c.Lecturer, c.Capacity).Scan(&c.ID, &c.Taken)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 = UNIQUE violation
			return domain.ErrDuplicateCode
		}
		// Membungkus error asli dari database[cite: 1]
		return fmt.Errorf("gagal insert course: %w", err)
	}
	return nil
}

func (r *courseRepository) GetByID(ctx context.Context, id int64) (domain.Course, error) {
	query := `SELECT id, code, name, description, lecturer, capacity, taken FROM courses WHERE id = $1`
	var c domain.Course
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Code, &c.Name, &c.Description, &c.Lecturer, &c.Capacity, &c.Taken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, domain.ErrCourseNotFound // Menerjemahkan ke sentinel error[cite: 1]
		}
		return c, fmt.Errorf("gagal get course: %w", err)
	}
	c.Remaining = c.Capacity - c.Taken
	return c, nil
}

// ReserveSeat menambah kursi terpakai secara atomik untuk mencegah Race Condition (BR-01)[cite: 1]
func (r *courseRepository) ReserveSeat(ctx context.Context, id int64) error {
	query := `UPDATE courses SET taken = taken + 1 WHERE id = $1 AND taken < capacity RETURNING id`
	var returnedID int64
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNoSeat // Kuota habis atau course tidak ada[cite: 1]
		}
		return fmt.Errorf("gagal reserve seat: %w", err)
	}
	return nil
}

// ReleaseSeat mengembalikan kursi jika enrollment dibatalkan (BR-03 & BR-06)[cite: 1]
func (r *courseRepository) ReleaseSeat(ctx context.Context, id int64) error {
	query := `UPDATE courses SET taken = taken - 1 WHERE id = $1 AND taken > 0`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal release seat: %w", err)
	}
	return nil
}