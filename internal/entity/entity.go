package entity

import (
	"github.com/xeynse/XeynseJar_analytics/internal/entity/home"
	"github.com/xeynse/XeynseJar_analytics/internal/util/error"
)

type Response struct {
	Data    interface{}  `json:"data"`
	Status  int          `json:"statusCode"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Error   *error.Error `json:"error"`
}

type AllJarStatusRequest struct {
	HomeID string `json:"homeID"`
}

type JarStatusRequest struct {
	HomeID string `json:"homeID"`
	JarID  string `json:"jarID"`
}

type ConsumptionRequest struct {
	HomeID  string `json:"homeID"`
	JarID   string `json:"jarID"`
	Context string `json:"context"`
}

type JarConsumption struct {
	TotalConsumption float64     `json:"total_consumption"`
	Data             interface{} `json:"data"`
}

type ConsumptionForday struct {
	Hour  int     `json:"hour"`
	Value float64 `json:"value"`
}

type ConsumptionForWeek struct {
	Date  string  `json:"date"`
	Day   string  `json:"day"`
	Value float64 `json:"value"`
}

type ConsumptionForMonth struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type JarReported struct {
	Nodes JarNodes `json:"nodes"`
}

type JarNodes struct {
	Node0 JarState `json:"0"`
}

type JarState struct {
	HomeID        string  `json:"home_id"`
	Ver           string  `json:"ver"`
	WeightStart   float64 `json:"w_start"`
	WeightDiff    float64 `json:"usedAmount"`
	WeightCurrent float64 `json:"w_curr"`
	ZAxisG        float64 `json:"zAxis_g"`
	Error         float64 `json:"jar_error"`
	Type          string  `json:"type"`
}

type HomeConfigurationResponse struct {
	HomeID        string             `json:"homeID"`
	Configuration home.Configuration `json:"configuration"`
}
