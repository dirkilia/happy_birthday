package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"happy_birthday/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer(log *slog.Logger) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	jwtKey := os.Getenv("JWT_SECRET")
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(log, NewServer.db, jwtKey),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
