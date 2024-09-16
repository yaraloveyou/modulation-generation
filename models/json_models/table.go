package jsonmodels

type Table struct {
	Name     string      `json:"name"`
	Modifier string      `json:"modifier"`
	Fields   interface{} `json:"fields"`
}
