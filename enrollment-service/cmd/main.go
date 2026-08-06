package main

import (
	"database/sql"
	"log"
	"os"

	"enrollment-service/internal/client"
	"enrollment-service/internal/handler"
	"enrollment-service/internal/repository"
	"enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5433/db_enrollment?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Gagal konek DB:", err)
	}
	defer db.Close()

	// Inisialisasi komponen
	courseClient := client.NewCourseClient()
	repo := repository.NewEnrollmentRepository(db)
	svc := service.NewEnrollmentService(repo, courseClient)
	h := handler.NewEnrollmentHandler(svc)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routing API
	r.POST("/enrollments", h.Enroll)
	r.GET("/enrollments/:id", h.Get)
	r.DELETE("/enrollments/:id", h.Cancel)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Port untuk Enrollment Service[cite: 1]
	}

	log.Printf("Enrollment Service berjalan di port %s", port)
	r.Run(":" + port)
}