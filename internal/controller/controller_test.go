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
