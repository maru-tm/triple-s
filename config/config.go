package config

type Config struct {
	Port       string
	StorageDir string
}

var GlobalConfig *Config

func NewConfig(port, storageDir string) *Config {
	GlobalConfig = &Config{
		Port:       port,
		StorageDir: storageDir,
	}
	return GlobalConfig
}

func GetStorageDir() string {
	return GlobalConfig.StorageDir
}
