package config

import (
	"log"
	"os"
	"regexp"

	"github.com/fvrvz/auth-service-go/dto"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var config *dto.Config

func resolveEnvPlaceholders(content string) string {
	re := regexp.MustCompile(`\$\{(\w+)(?::-(.*?))?\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		key := parts[1]
		defaultVal := parts[2]
		if val := os.Getenv(key); val != "" {
			return val
		}
		return defaultVal
	})
}

func Init(yamlPath string, envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	raw, err := os.ReadFile(yamlPath)

	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	parsed := resolveEnvPlaceholders(string(raw))

	config = &dto.Config{}

	if err := yaml.Unmarshal([]byte(parsed), config); err != nil {
		log.Fatalf("Failed to serialize config: %v", err)
	}

	log.Println("Config loaded successfully")
}

func GetConfig() *dto.Config {
	return config
}
