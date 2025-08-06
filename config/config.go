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
	re := regexp.MustCompile(`\$\{(\w+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		return os.Getenv(key)
	})
}

func Init(yamlPath string, envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	raw, err := os.ReadFile(yamlPath)

	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
		return
	}

	parsed := resolveEnvPlaceholders(string(raw))

	config = &dto.Config{}

	if err := yaml.Unmarshal([]byte(parsed), config); err != nil {
		log.Fatalf("Failed to serialize config: %v", err)
		return
	}

	log.Println("Config loaded successfully")
}

func GetConfig() *dto.Config {
	return config
}
