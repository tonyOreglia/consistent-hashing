package consistent_hash

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server abstraction.
type Server struct {
	r      *gin.Engine
	vNodes []virtualNode
	nodes  map[string]string // node-id to node URL map
}

// NewServer instantiates a new HTTP Server.
func NewServer() (s *Server) {
	s = &Server{
		r:     gin.Default(),
		nodes: make(map[string]string),
	}

	s.r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// POST /nodes -> node ID string
	s.r.POST("/nodes", s.AddNode)
	s.r.GET("/nodes/count", s.NodeCount)
	s.r.DELETE("/nodes/:nodeId", s.DeleteNode)

	// s.r.GET("/exists/:word", s.WordExists)
	// s.r.POST("/add", s.Add)
	// s.r.GET("/matches/:prefix", s.Matches)

	return
}

// Run the HTTP server. It will block until a fatal error is encountered.
func (s *Server) Run(addr string) error {
	return s.r.Run(addr)
}
