package room

import (
	deviceentity "github.com/xeynse/XeynseJar_analytics/internal/entity/device"
)

type Room struct {
	ID      string                         `json:"id"`
	Name    string                         `json:"name"`
	Devices map[string]deviceentity.Device `json:"devices"`
}
