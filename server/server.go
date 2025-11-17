package server

import (
	"net/http"

	"consistent_hash/controller"

	"github.com/gin-gonic/gin"
)

// Server abstraction.
type Server struct {
	r          *gin.Engine
	controller *controller.Controller
}

// NewServer instantiates a new HTTP Server.
func NewServer() (s *Server) {
	s = &Server{
		r:          gin.Default(),
		controller: controller.NewController(),
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
