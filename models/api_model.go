package models

type ParamType string

// Type constants for parameter validation
const (
	TypeString    ParamType = "String"
	TypeInt       ParamType = "Int"
	TypeBoolean   ParamType = "Boolean"
	TypeList      ParamType = "List"
	TypeDate      ParamType = "Date"
	TypeEmail     ParamType = "Email"
	TypeUUID      ParamType = "UUID"
	TypeAuthToken ParamType = "Auth-Token"
)

type Parameter struct {
	Name     string      `json:"name"`
	Types    []ParamType `json:"types"`
	Required bool        `json:"required"`
}

type APIModel struct {
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	QueryParams []*Parameter `json:"query_params"`
	Headers     []*Parameter `json:"headers"`
	Body        []*Parameter `json:"body"`
}
