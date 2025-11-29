package models

import "fmt"

// Type constants for parameter validation
const (
	TypeString    = "String"
	TypeInt       = "Int"
	TypeBoolean   = "Boolean"
	TypeList      = "List"
	TypeDate      = "Date"
	TypeEmail     = "Email"
	TypeUUID      = "UUID"
	TypeAuthToken = "Auth-Token"
)

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
func Key(path, method string) string {
	return fmt.Sprintf("%s:%s", path, method)
}
