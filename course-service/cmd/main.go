package main

import (
	"database/sql"
	"log"
	"os"

	"course-service/internal/handler"
	"course-service/internal/repository"
	"course-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Baca kredensial dari environment variable (bukan hardcode)[cite: 1]
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback untuk testing lokal jika dijalankan tanpa docker
		dbURL = "postgres://user:password@localhost:5432/db_course?sslmode=disable"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Gagal konek DB:", err)
	}
	defer db.Close()

	// Dependency Injection[cite: 1]
	repo := repository.NewCourseRepository(db)
	svc := service.NewCourseService(repo)
	h := handler.NewCourseHandler(svc)

	r := gin.Default()
	
	// Middleware logging dan recovery bawaan Gin[cite: 1]
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routing API
	r.POST("/courses", h.Create)
	r.GET("/courses/:id", h.Get)
	r.GET("/courses/:id/availability", h.GetAvailability)

	r.POST("/courses/:id/seats/reserve", h.Reserve)
	r.POST("/courses/:id/seats/release", h.Release)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Standar port Course Service[cite: 1]
	}

	log.Printf("Course Service berjalan di port %s", port)
	r.Run(":" + port)
}