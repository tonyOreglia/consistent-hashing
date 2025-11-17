package node

import "consistent_hash/internal/redis"

type VirtualNode struct {
	HashValue uint64
	NodeId    string // 8 byte hash id derived from node URL
	Client    *redis.Client
}
