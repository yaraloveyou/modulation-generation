package config

type Config struct {
	IsLombok     bool
	OutputDir    string
	FolderStruct string
}

func GetConfig() *Config {
	return &Config{
		IsLombok:     true,
		OutputDir:    "../../../generated_sources",
		FolderStruct: "../../config/folder_struct.json",
	}
}
