package server

import (
	"log"
	"net/http"

	"consistent_hash/internal/controller"
	"consistent_hash/internal/server/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Server abstraction.
type Server struct {
	r          *gin.Engine
	controller *controller.Controller
	config     *config.Config
}

// NewServer instantiates a new HTTP Server.
func NewServer() (s *Server) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration from environment variables
	config := config.NewConfig()

	s = &Server{
		r:          gin.Default(),
		controller: controller.NewController(config),
		config:     config,
	}

	s.r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s.r.POST("/nodes", s.AddNode)
	s.r.GET("/nodes/count", s.NodeCount)
	s.r.DELETE("/nodes/:nodeId", s.DeleteNode)
	s.r.GET("/nodes/:key", s.GetNodes)

	return
}

// Run the HTTP server. It will block until a fatal error is encountered.
func (s *Server) Run(addr string) error {
	return s.r.Run(addr)
}
