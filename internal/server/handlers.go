package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddNodeRequest struct {
	Url string
}

type DeleteNodeResponse struct {
	NodeUrl string `json:"nodeUrl"`
}

type NodeResponse struct {
	NodeId string `json:"nodeId"`
}

func (s *Server) AddNode(c *gin.Context) {
	var req AddNodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiError{Err: err, Desc: "unable to parse request"})
		return
	}
	nodeUrl := req.Url

	log.Println("Adding node at ", nodeUrl)

	nodeId, err := s.controller.AddNode(nodeUrl, s.redisFactory)

	if err != nil {
		c.JSON(http.StatusNotFound, ApiError{Err: err, Desc: "problem adding node"})
		return
	}

	c.JSON(http.StatusOK, NodeResponse{NodeId: nodeId})
}

func (s *Server) NodeCount(c *gin.Context) {
	resp := s.controller.NodeCount()

	c.JSON(http.StatusOK, resp)
}

func (s *Server) DeleteNode(c *gin.Context) {
	nodeId := c.Param("nodeId")

	url, err := s.controller.DeleteNode(nodeId)

	if err != nil {
		c.JSON(http.StatusNotFound, ApiError{Err: err, Desc: "problem deleting node"})
		return
	}

	c.JSON(http.StatusOK, DeleteNodeResponse{NodeUrl: url})
}

func (s *Server) GetNodes(c *gin.Context) {
	key := c.Param("key")
	resp := s.controller.GetNodes(key)

	c.JSON(http.StatusOK, resp)
}

type UpsertKeyValueRequest struct {
	Content string `json:"content"`
}

func (s *Server) StoreValue(c *gin.Context) {
	key := c.Param("key")
	var req UpsertKeyValueRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiError{Err: err, Desc: "unable to parse request"})
		return
	}
	err := s.controller.StoreValue(key, req.Content)
	if err != nil {
		c.JSON(http.StatusNotFound, ApiError{Err: err, Desc: "problem storing value"})
		return
	}
}

type GetValueResponse struct {
	Content string `json:"content"`
}

func (s *Server) GetValue(c *gin.Context) {
	key := c.Param("key")
	content, err := s.controller.GetValue(key)
	if err != nil {
		c.JSON(http.StatusNotFound, ApiError{Err: err, Desc: "problem retrieving value"})
		return
	}
	c.JSON(http.StatusOK, GetValueResponse{Content: content})

}
