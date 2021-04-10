package entity

import "github.com/xeynse/XeynseJar_analytics/internal/util/error"

type Response struct {
	Data    interface{}  `json:"data"`
	Status  int          `json:"statusCode"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Error   *error.Error `json:"error"`
}

type JarState struct {
	HomeID        string  `json:"home_id"`
	WeightStart   float64 `json:"w_start"`
	WeightDiff    float64 `json:"w_diff"`
	WeightCurrent float64 `json:"w_curr"`
	ZAxisG        float64 `json:"zAxis_g"`
	Error         float64 `json:"jar_error"`
	Type          string  `json:"type"`
}
