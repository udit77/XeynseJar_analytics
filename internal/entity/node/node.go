package node

type Node struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
	Properties string             `json:"properties"`
}

type ModeNode struct{
	ID     string                 `json:"id"`
	State  map[string]interface{} `json:"state"`
}