package models

type RequestParam struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type Request struct {
	Path        string          `json:"path"`
	Method      string          `json:"method"`
	QueryParams []*RequestParam `json:"query_params"`
	Headers     []*RequestParam `json:"headers"`
	Body        []*RequestParam `json:"body"`
}
