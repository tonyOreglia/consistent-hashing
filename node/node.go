package node

type VirtualNode struct {
	HashValue uint64
	NodeId    string // 8 byte hash id derived from node URL
}
type Node struct {
	Url    string `json:"url"`
	NodeId string `json:"nodeId"`
}
