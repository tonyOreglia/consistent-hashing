package controller

import (
	"consistent_hash/internal/hash"
	"consistent_hash/internal/node"
	"consistent_hash/internal/redis"
	"consistent_hash/internal/server/config"

	"fmt"
	"log"
	"sort"
)

type RedisFactory interface {
	New(nodeUrl string) (*redis.Client, error)
}

type Controller struct {
	vNodes        []node.VirtualNode
	nodesUrlsById map[string]string // node-id to node URL map
	config        *config.Config
}

func NewController(config *config.Config) (c *Controller) {
	c = &Controller{
		nodesUrlsById: make(map[string]string),
		config:        config,
	}

	return c
}

func (c *Controller) AddNode(nodeUrl string, redis RedisFactory) (string, error) {
	log.Println("Adding node at ", nodeUrl)

	client, err := redis.New(nodeUrl)

	if err != nil {
		return "", fmt.Errorf("unable to connect with node at %s: %v", nodeUrl, err)
	}

	nodeId := hash.HashId(nodeUrl)

	if _, ok := c.nodesUrlsById[nodeId]; ok {
		return "", fmt.Errorf("node already exists with id %s", nodeId)
	}

	c.nodesUrlsById[nodeId] = nodeUrl
	log.Println("Derived Node ID ", nodeId)

	// create 100 virtual nodes for this server
	for i := 0; i < 10; i++ {
		hash := hash.HashKey(fmt.Sprintf("%s_%d", nodeUrl, i))
		c.vNodes = append(c.vNodes, node.VirtualNode{HashValue: hash, NodeId: nodeId, Client: client})
	}

	sort.Slice(c.vNodes, func(i, j int) bool {
		return c.vNodes[i].HashValue < c.vNodes[j].HashValue
	})

	return nodeId, nil
}

type NodeCountResponse struct {
	NodeCount        int `json:"nodeCount"`
	VirtualNodeCount int `json:"virtualNodeCount"`
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

func (c *Controller) GetNodes(key string) []node.VirtualNode {
	kHash := hash.HashKey(key)

	targetNode := c.vNodes[0]
	res := []node.VirtualNode{}

	for i, value := range c.vNodes {
		log.Printf("checking virtual node at index %d", i)
		if kHash > value.HashValue {
			continue
		}

		targetNode = c.vNodes[i]
		res = []node.VirtualNode{{NodeId: targetNode.NodeId}}

		j := i

		for len(res) < c.config.Redundancy && len(res) < len(c.nodesUrlsById) {
			j++
			if j == len(c.vNodes) {
				j = 0
			}
			newNode := node.VirtualNode{NodeId: c.vNodes[j].NodeId}
			if !isUnique(newNode, res) {
				j++
				continue
			}
			res = append(res, newNode)
		}
		break

	}

	if len(res) == 0 {
		res = []node.VirtualNode{{NodeId: c.vNodes[0].NodeId}}
	}

	return res
}

func isUnique(newNode node.VirtualNode, nodes []node.VirtualNode) bool {
	for _, v := range nodes {
		if v == newNode {
			return false
		}
	}
	return true
}

func (c *Controller) StoreValue(key string, content string) error {
	targetNodes := c.findNodes(key)

	for _, node := range targetNodes {
		log.Printf("storing content for key %s on node %s", key, node.NodeId)
		err := node.Client.Set(key, content)
		if err != nil {
			return fmt.Errorf("unable to set value on node %s: %v", node.NodeId, err)
		}
	}
	return nil
}

func (c *Controller) findNodes(key string) []node.VirtualNode {
	kHash := hash.HashKey(key)
	targetNodes := []node.VirtualNode{}

	for i, value := range c.vNodes {
		if kHash > value.HashValue {
			continue
		}

		targetNode := c.vNodes[i]
		targetNodes = append(targetNodes, targetNode)

		j := i

		for len(targetNodes) < c.config.Redundancy && len(targetNodes) < len(c.nodesUrlsById) {
			j++
			if j == len(c.vNodes) {
				j = 0
			}
			if !isUnique(c.vNodes[j], targetNodes) {
				j++
				continue
			}
			targetNodes = append(targetNodes, c.vNodes[j])
		}
		break

	}

	if len(targetNodes) == 0 {
		targetNodes = append(targetNodes, c.vNodes[0])
	}
	return targetNodes
}

func (c *Controller) GetValue(key string) (string, error) {
	targetNodes := c.findNodes(key)
	contents := []string{}

	for _, node := range targetNodes {
		content, err := node.Client.Get(key)
		if err != nil {
			return "", fmt.Errorf("failed retrieving value from node %s: %v", node.NodeId, err)
		}
		log.Printf("got value %s for key %s from node %s", content, key, node.NodeId)
		contents = append(contents, content)
	}
	for i, v := range contents {
		if v != contents[0] {
			return "", fmt.Errorf("content mismatch %s from node %s != %s from node %s", contents[0], targetNodes[0].NodeId, contents[i], targetNodes[i].NodeId)
		}
	}
	return contents[0], nil
}
