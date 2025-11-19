package controller

import (
	"errors"
	"fmt"
	"testing"

	"consistent_hash/internal/redis"
	"consistent_hash/internal/server/config"
)

type mockedRedisFactory struct{}

func (m mockedRedisFactory) New(url string) (*redis.Client, error) {
	c := &redis.Client{}
	return c, nil
}

func TestDeleteNode(t *testing.T) {
	cfg := config.NewConfig()
	mockedRedis := &mockedRedisFactory{}
	cntr := NewController(cfg)
	url1 := "node-1"
	node1Id, err := cntr.AddNode(url1, mockedRedis)
	if err != nil {
		t.Fatalf("error adding node1: %v", err)
	}
	cntr.AddNode("node-2", mockedRedis)
	cntr.AddNode("node-3", mockedRedis)
	url4 := "node-4"
	node4Id, err := cntr.AddNode(url4, mockedRedis)
	if err != nil {
		t.Fatalf("error adding node1: %v", err)
	}
	// node5Id := cntr.AddNode("node-5", mockedRedis)

	deletedUrl, err := cntr.DeleteNode(node4Id)
	if err != nil {
		t.Fatalf("error deleting node1: %v", err)
	}
	if deletedUrl != url4 {
		t.Fatalf("deleted URL [%s] not equal to url4 [%s]", deletedUrl, url4)
	}

	_, err = cntr.DeleteNode("missing-node")
	if err == nil {
		t.Fatalf("expected error deleting missing-node")
	}
	if !errors.Is(err, ErrNodeDoesNotExist) {
		t.Fatalf("expecting ErrNodeDoesNotExist type error, got %v", err)
	}

	deletedUrl, err = cntr.DeleteNode(node1Id)
	if err != nil {
		t.Fatalf("error deleting node1: %v", err)
	}
	if deletedUrl != url1 {
		t.Fatalf("deleted URL [%s] not equal to url1 [%s]", deletedUrl, url1)
	}

	// cannot delete same node twice
	_, err = cntr.DeleteNode(node1Id)
	if err == nil {
		t.Fatalf("expected error deleting missing-node")
	}
	if !errors.Is(err, ErrNodeDoesNotExist) {
		t.Fatalf("expecting ErrNodeDoesNotExist type error, got %v", err)
	}
}

func TestAddNode(t *testing.T) {
	cfg := config.NewConfig()
	mockedRedis := &mockedRedisFactory{}

	type expectedResult struct {
		nodeId       string
		addNodeError error
	}
	tests := []struct {
		urls           []string
		expectedResult []expectedResult
		name           string
	}{
		{name: "add single node valid url", urls: []string{"http://test:1234"}, expectedResult: []expectedResult{{nodeId: "b8db6c8ad1cd7f36"}}},
		{
			name: "add multiple nodes valid url", urls: []string{
				"http://test:1234", "http://test:5678", "http://test:9098",
			},
			expectedResult: []expectedResult{{nodeId: "b8db6c8ad1cd7f36"}, {nodeId: "b87fd3cb5092cca8"}, {nodeId: "4496693eda7e228a"}},
		},
		{
			name: "cannot add same node twice", urls: []string{
				"http://test:1234", "http://test:1234",
			},
			expectedResult: []expectedResult{{nodeId: "b8db6c8ad1cd7f36"}, {addNodeError: ErrNodeAlreadyExists}},
		},
		{
			name: "cannot add same node twice v2", urls: []string{
				"same-node", "unique-node", "same-node",
			},
			expectedResult: []expectedResult{{nodeId: "0db8ab6e506f3b10"}, {nodeId: "3a5695e306800e7f"}, {addNodeError: ErrNodeAlreadyExists}},
		},
	}

	for _, test := range tests {
		fmt.Printf("\n\nRunning test: %s\n", test.name)
		cntr := NewController(cfg)
		for i, url := range test.urls {
			nodeId, err := cntr.AddNode(url, mockedRedis)
			expectedResult := test.expectedResult[i]
			if err != nil && expectedResult.addNodeError != nil {
				if errors.Is(err, expectedResult.addNodeError) {
					continue
				}
				t.Fatalf("expected error %v, got error %v", expectedResult.addNodeError, err)
			}
			if err != nil {
				t.Fatalf("[%s] Error adding node: %v", test.name, err)
			}
			if expectedResult.addNodeError != nil {
				t.Fatalf("expected error %v, got no error", expectedResult.addNodeError)
			}

			expectedNodeId := expectedResult.nodeId
			if nodeId != expectedNodeId {
				t.Fatalf("[%s] Actual nodeId [%s] does not match expected node Id [%s]", test.name, nodeId, expectedNodeId)
			}
		}
	}
}

func TestFindNodes(t *testing.T) {
	cfg := &config.Config{
		ReadQuorum:  1,
		WriteQuorum: 1,
		Redundancy:  1,
	}
	mockedRedis := &mockedRedisFactory{}
	cntr := NewController(cfg)
	url1 := "node-1"
	_, err := cntr.AddNode(url1, mockedRedis)
	if err != nil {
		t.Fatalf("error adding node1: %v", err)
	}
	n2Id, _ := cntr.AddNode("node-2", mockedRedis)
	node3Id, _ := cntr.AddNode("node-3", mockedRedis)
	url4 := "node-4"
	_, err = cntr.AddNode(url4, mockedRedis)
	if err != nil {
		t.Fatalf("error adding node1: %v", err)
	}
	res := cntr.findNodes("key-1")
	if len(res) != 1 {
		t.Fatalf("expected single result with redundancy 1")
	}
	if res[0].NodeId != node3Id {
		t.Fatalf("expected node ID %s, but got %s", node3Id, res[0].NodeId)
	}
	cntr.config = &config.Config{
		ReadQuorum:  1,
		WriteQuorum: 1,
		Redundancy:  2,
	}
	res = cntr.findNodes("key-1")
	if len(res) != 2 {
		t.Fatalf("expected single result with redundancy 1")
	}
	if res[0].NodeId != node3Id {
		t.Fatalf("expected node ID %s, but got %s", node3Id, res[0].NodeId)
	}
	if res[1].NodeId != n2Id {
		t.Fatalf("expected node ID %s, but got %s", n2Id, res[1].NodeId)
	}
}
