// Package server contains HTTP handlers and server setup for the KV store API.
package server

import (
	"kv-store-wal/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler represents the HTTP server that handles API requests.
type Handler struct {
	router *gin.Engine
}

// NewHandler creates and initializes a new server instance.
func NewHandler() *Handler {
	router := gin.Default()

	s := &Handler{
		router: router,
	}

	s.setupRoutes()
	return s
}

// Handler returns the HTTP handler for the server.
func (s *Handler) Handler() http.Handler {
	return s.router
}

// setupRoutes configures all the routes for the server.
func (s *Handler) setupRoutes() {
	subRouter := s.router.Group("/v1")
	{
		subRouter.GET("/:key", s.Get)
		subRouter.POST("/write", s.Set)
	}
}

// Get handles requests to retrieve a value by key.
func (s *Handler) Get(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.String(http.StatusBadRequest, "Bad Request = key required")
		return
	}

	val, err := service.Get(c.Request.Context(), key)
	if err != nil {
		c.String(http.StatusNotFound, "Key not found")
		return
	}

	type Resp struct {
		Value string `json:"value"`
	}

	response := Resp{
		Value: val,
	}

	c.JSON(http.StatusOK, response)
}

// Set handles requests to store a key-value pair.
func (s *Handler) Set(c *gin.Context) {
	req := service.KVPair{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Key == "" || req.Value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "both key and value must be set",
		})
		return
	}

	err = service.Set(c.Request.Context(), req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
