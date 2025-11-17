package controller

import (
	"consistent_hash/hash"
	"consistent_hash/node"
	"fmt"
	"log"
	"sort"
)

type Controller struct {
	vNodes        []node.VirtualNode
	nodesUrlsById map[string]string // node-id to node URL map
}

func NewController() (c *Controller) {
	c = &Controller{
		nodesUrlsById: make(map[string]string),
	}

	return c
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

func (c *Controller) AddNode(nodeUrl string) (string, error) {
	log.Println("Adding node at ", nodeUrl)

	nodeId := hash.HashId(nodeUrl)

	if _, ok := c.nodesUrlsById[nodeId]; ok {
		return "", fmt.Errorf("node already exists with id %s", nodeId)
	}

	c.nodesUrlsById[nodeId] = nodeUrl
	log.Println("Derived Node ID ", nodeId)

	// create 100 virtual nodes for this server
	for i := 0; i < 10; i++ {
		hash := hash.HashKey(fmt.Sprintf("%s_%d", nodeUrl, i))
		c.vNodes = append(c.vNodes, node.VirtualNode{HashValue: hash, NodeId: nodeId})
	}

	sort.Slice(c.vNodes, func(i, j int) bool {
		return c.vNodes[i].HashValue < c.vNodes[j].HashValue
	})

	return nodeId, nil
}

func (c *Controller) NodeCount() NodeCountResponse {
	nodeCount := len(c.nodesUrlsById)
	vNodeCount := len(c.vNodes)

	return NodeCountResponse{NodeCount: nodeCount, VirtualNodeCount: vNodeCount}
}

func (c *Controller) DeleteNode(nodeId string) (string, error) {
	url, ok := c.nodesUrlsById[nodeId]

	if !ok {
		return "", fmt.Errorf("node does not exist")
	}

	delete(c.nodesUrlsById, nodeId)

	var newVirtualNodes []node.VirtualNode

	for _, vNode := range c.vNodes {
		if vNode.NodeId != nodeId {
			newVirtualNodes = append(newVirtualNodes, vNode)
		}
	}
	c.vNodes = newVirtualNodes

	return url, nil
}

func (c *Controller) GetNodes(key string) []node.Node {
	kHash := hash.HashKey(key)

	targetNode := c.vNodes[0]
	for index, value := range c.vNodes {
		log.Println("Checking virtual node at index %s", index)
		if kHash < value.HashValue {
			targetNode = c.vNodes[index]
			break
		}
	}
	return []node.Node{{NodeId: targetNode.NodeId, Url: c.nodesUrlsById[targetNode.NodeId]}}
}
