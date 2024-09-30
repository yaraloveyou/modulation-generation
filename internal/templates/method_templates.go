package templates

const (
	GETTER_TEMPLATE = `{{.Modifier}} {{.Return.DataType}} {{.Name}}({{StringsJoin .ExternalVariables ", "}}) {
		return {{.Return.Name}}
	}`

	SETTER_TEMPLATE = `{{.Modifier}} {{.Return.DataType}} {{.Name}}({{with index .Variables 0}}{{.DataType}} {{.Name}}{{end}}) {
		this.{{with index .ExternalVariables 0}}{{.Name}}{{end}} = {{with index .Variables 0}}{{.Name}}{{end}};
	}`

	CONSTRUCTOR_TEMPLATE = `{{.Modifier}} {{.Name}}({{StringsJoin .ExternalVariables ", "}}) {
		{{VarsJoin .ExternalVariables "\n\t\t"}}	
	}`

	TOSTRING_TEMPLATE = `{{.Modifier}} {{.Return.DataType}} {{.Name}}() {
		return String.format("{{.ClassName}}[{{StringsJoin .ExternalVariables ", "}}]", {{VarsJoin .ExternalVariables ", "}})	
	}`
)
