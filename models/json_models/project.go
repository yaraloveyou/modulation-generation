package jsonmodels

type Project struct {
	Name   string
	Lombok bool
	Tables []Table
}
