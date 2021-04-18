package device

import nodeentity "github.com/xeynse/XeynseJar_analytics/internal/entity/node"

type Device struct {
	ID    string                     `json:"id"`
	Nodes map[string]nodeentity.Node `json:"nodes"`
}
