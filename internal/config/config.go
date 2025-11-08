package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenAI struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"openai"`
}

func LoadConfig(filename string) Config {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file %s: %v", filename, err))
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config file %s: %v", filename, err))
	}

	return config
}
