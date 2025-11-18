package controller

import (
	"consistent_hash/internal/redis"
	"consistent_hash/internal/server/config"
	"testing"
)

type mockedRedisFactory struct{}

func (m mockedRedisFactory) New(url string) (*redis.Client, error) {
	c := &redis.Client{}
	return c, nil
}

func TestAddNode(t *testing.T) {
	cfg := config.NewConfig()
	cntr := NewController(cfg)
	mockedRedis := &mockedRedisFactory{}

	tests := []struct {
		url            string
		expectedNodeId string
	}{
		{url: "http://test:1234", expectedNodeId: "b8db6c8ad1cd7f36"},
	}

	for _, test := range tests {

		nodeId, err := cntr.AddNode(test.url, mockedRedis)
		if err != nil {
			t.Errorf("Error adding node: %v", err)
		}

		if nodeId != test.expectedNodeId {
			t.Errorf("Actual nodeId [%s] does not match expected node Id [%s]", nodeId, test.expectedNodeId)
		}
	}

}
