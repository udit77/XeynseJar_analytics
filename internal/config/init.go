package config

import "github.com/xeynse/XeynseJar_analytics/internal/resource/file"

var config *Config

func Init(file file.Resource) (*Config, error) {
	configFileName := "xeynseJar_analytics.json"
	err := file.ReadJSONFile(configFileName, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func GetConfig() *Config {
	return config
}
