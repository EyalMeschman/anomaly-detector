package models

type FieldAnomaly struct {
	Field         string `json:"field"`
	ParameterName string `json:"parameter_name"`
	Reason        string `json:"reason"`
}

type ValidationResult struct {
	Valid     bool            `json:"valid"`
	Anomalies []*FieldAnomaly `json:"anomalies,omitempty"`
}
