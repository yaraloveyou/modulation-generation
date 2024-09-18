package jsonmodels

type FileSystem struct {
	Pkg     string `json:"package"`
	Folders map[string]string
}
