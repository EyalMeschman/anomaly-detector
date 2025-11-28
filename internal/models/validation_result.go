package models

type FieldAnomaly struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Reason   string `json:"reason"`
}

type ValidationResult struct {
	IsAnomalous     bool           `json:"is_anomalous"`
	AnomalousFields []FieldAnomaly `json:"anomalous_fields,omitempty"`
}
