package generator

func BuildProject() error {
	return generateProject("../../example/tables.json")
}

func generateProject(path string) error {
	if err := GenerateClasses(path); err != nil {
		return err
	}
	return GenerateFiles()
}
