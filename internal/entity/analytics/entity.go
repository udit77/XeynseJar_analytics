package entity

import "time"

type CurrentWeightStatus struct {
	JarID         string  `db:"jar_id" json:"id"`
	WeightCurrent float64 `db:"weight_current" json:"quantity"`
}

type Consumption struct {
	JarID      string    `db:"jar_id"`
	WeightDiff float64   `db:"weight_diff"`
	UpdateTime time.Time `db:"update_time"`
}
