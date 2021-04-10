package entity

type CurrentWeightStatus struct {
	JarID         string  `db:"jar_id" json:"id"`
	WeightCurrent float64 `db:"weight_current" json:"quantity"`
}
