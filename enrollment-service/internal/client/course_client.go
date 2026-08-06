package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"enrollment-service/internal/domain"
)

type CourseClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCourseClient() *CourseClient {
	url := os.Getenv("COURSE_SERVICE_URL")
	if url == "" {
		url = "http://course-service:8081" // DNS internal dari Docker Compose[cite: 1]
	}
	return &CourseClient{
		baseURL: url,
		// Timeout bawaan HTTP client (lapisan pertama perlindungan)
		httpClient: &http.Client{Timeout: 4 * time.Second},
	}
}

// GetCourse memanggil GET /courses/{id}
func (c *CourseClient) GetCourse(ctx context.Context, courseID int64) (string, error) {
	// BR-05: Wajib memakai context dengan timeout maksimal 3 detik[cite: 1]
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/courses/%d", c.baseURL, courseID), nil)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Mengubah error timeout menjadi error domain (BR-05)[cite: 1]
		return "", domain.ErrUpstreamUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", domain.ErrCourseNotFound // BR-04[cite: 1]
	}
	if resp.StatusCode != http.StatusOK {
		return "", domain.ErrUpstreamUnavailable
	}

	// Parsing struktur response {"success": true, "data": {...}}
	var body struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("gagal decode response: %w", err)
	}

	return body.Data.Name, nil
}

// Reserve memanggil POST /courses/{id}/seats/reserve
func (c *CourseClient) Reserve(ctx context.Context, courseID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // BR-05[cite: 1]
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/courses/%d/seats/reserve", c.baseURL, courseID), nil)
	if err != nil {
		return fmt.Errorf("gagal membuat request reserve: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ErrUpstreamUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return domain.ErrNoSeat // BR-01: Kuota penuh[cite: 1]
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrCourseNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.ErrUpstreamUnavailable
	}

	return nil
}

// Release memanggil POST /courses/{id}/seats/release (Kompensasi BR-03 & BR-06)[cite: 1]
func (c *CourseClient) Release(ctx context.Context, courseID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/courses/%d/seats/release", c.baseURL, courseID), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("release gagal dengan status %d", resp.StatusCode)
	}
	return nil
}