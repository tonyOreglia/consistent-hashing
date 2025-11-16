package consistent_hash

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type virtualNode struct {
	hashValue uint64
	nodeId    string // 8 byte hash id derived from node URL
}

type AddNodeRequest struct {
	Url string
}

type DeleteNodeResponse struct {
	NodeUrl string `json:"nodeUrl"`
}

type NodeResponse struct {
	NodeId string `json:"nodeId"`
}

type NodeCountResponse struct {
	NodeCount        int `json:"nodeCount"`
	VirtualNodeCount int `json:"virtualNodeCount"`
}

func (s *Server) AddNode(c *gin.Context) {
	// extract url from body

	var req AddNodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiError{Err: err, Desc: "unable to parse request"})
		return
	}
	nodeUrl := req.Url

	log.Println("Adding node at ", nodeUrl)

	nodeId := hashId(nodeUrl)

	if _, ok := s.nodes[nodeId]; ok {
		c.JSON(http.StatusBadRequest, ApiError{Err: fmt.Errorf("node already exists with id %s", nodeId), Desc: "node already exists"})
		return
	}

	s.nodes[nodeId] = nodeUrl
	log.Println("Derived Node ID ", nodeId)

	// create 100 virtual nodes for this server
	for i := 0; i < 10; i++ {
		hash := hashKey(fmt.Sprintf("%s_%d", nodeUrl, i))
		s.vNodes = append(s.vNodes, virtualNode{hashValue: hash, nodeId: nodeId})
	}

	sort.Slice(s.vNodes, func(i, j int) bool {
		return s.vNodes[i].hashValue < s.vNodes[j].hashValue
	})

	c.JSON(http.StatusOK, NodeResponse{NodeId: nodeId})
}

func (s *Server) NodeCount(c *gin.Context) {
	// count nodes
	nodeCount := len(s.nodes)
	vNodeCount := len(s.vNodes)

	c.JSON(http.StatusOK, NodeCountResponse{NodeCount: nodeCount, VirtualNodeCount: vNodeCount})
}

func (s *Server) DeleteNode(c *gin.Context) {
	nodeId := c.Param("nodeId")

	url, ok := s.nodes[nodeId]

	if !ok {
		c.JSON(http.StatusNotFound, ApiError{Err: fmt.Errorf("node does not exist"), Desc: "node not found"})
		return
	}

	delete(s.nodes, nodeId)

	var newVirtualNodes []virtualNode

	for _, vNode := range s.vNodes {
		if vNode.nodeId != nodeId {
			newVirtualNodes = append(newVirtualNodes, vNode)
		}
	}
	s.vNodes = newVirtualNodes

	c.JSON(http.StatusOK, DeleteNodeResponse{NodeUrl: url})
}

type Node struct {
	Url    string `json:"url"`
	NodeId string `json:"nodeId"`
}
type GetNodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func (s *Server) GetNodes(c *gin.Context) {
	key := c.Param("key")
	kHash := hashKey(key)

	targetNode := s.vNodes[0]
	for index, value := range s.vNodes {
		log.Println("Checking virtual node at index %s", index)
		if kHash < value.hashValue {
			targetNode = s.vNodes[index]
			break
		}
	}
	resp := GetNodesResponse{
		Nodes: []Node{{NodeId: targetNode.nodeId, Url: s.nodes[targetNode.nodeId]}},
	}
	c.JSON(http.StatusOK, resp)
}
