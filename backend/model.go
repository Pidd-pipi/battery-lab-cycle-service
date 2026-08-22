package main

type Cell struct {
	ID          string  `json:"id"`
	Chemistry   string  `json:"chemistry"`
	Cycle       int     `json:"cycle"`
	CapacityMah float64 `json:"capacity_mah"`
	VoltageV    float64 `json:"voltage_v"`
	Status      string  `json:"status"`
}
