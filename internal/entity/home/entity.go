package home

import (
	roomentity "github.com/xeynse/XeynseJar_analytics/internal/entity/room"
)

type Configuration struct {
	HomeID string                     `json:"homeID"`
	Rooms  map[string]roomentity.Room `json:"rooms"`
}
