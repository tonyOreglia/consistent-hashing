package node

type VirtualNode struct {
	HashValue uint64
	NodeId    string // 8 byte hash id derived from node URL
}
