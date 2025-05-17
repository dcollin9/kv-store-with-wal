package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler represents the HTTP server that handles API requests
type Handler struct {
	router *gin.Engine
}

// NewHandler creates and initializes a new server instance
func NewHandler() *Handler {
	router := gin.Default()

	s := &Handler{
		router: router,
	}

	s.setupRoutes()
	return s
}

// Handler returns the HTTP handler for the server
func (s *Handler) Handler() http.Handler {
	return s.router
}

// setupRoutes configures all the routes for the server
func (s *Handler) setupRoutes() {
	// Root route
	s.router.GET("/", s.handleIndex)

	// API routes
	api := s.router.Group("/api")
	{
		api.GET("/hello", s.handleHello)
	}
}

// handleIndex handles the root path
func (s *Handler) handleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "KV Store with WAL API is running",
	})
}

// handleHello is an example endpoint that accepts HTTP requests
func (s *Handler) handleHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Hello, world!",
		},
	})
}
