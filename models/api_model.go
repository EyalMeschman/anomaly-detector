package models

import "fmt"

type Parameter struct {
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	Required bool     `json:"required"`
}

type APIModel struct {
	Path        string      `json:"path"`
	Method      string      `json:"method"`
	QueryParams []Parameter `json:"query_params"`
	Headers     []Parameter `json:"headers"`
	Body        []Parameter `json:"body"`
}

// Key returns a unique identifier for this API model
func (m *APIModel) Key() string {
	return fmt.Sprintf("%s:%s", m.Path, m.Method)
}
