package jsonmodels

type Table struct {
	Name   string      `json:"name"`
	Fields interface{} `json:"fields"`
}
