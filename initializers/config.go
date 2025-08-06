package initializers

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Port int `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	SSLMode  string `yaml:"sslmode"`
}

type JWT struct {
	ExpiryHours int    `yaml:"expiry_hours"`
	JWTSecret   string `yaml:"jwt_secret"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`
}

func resolveEnvPlaceholders(content string) string {
	re := regexp.MustCompile(`\$\{(\w+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		return os.Getenv(key)
	})
}

func LoadConfig(yamlPath string, envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(yamlPath)

	if err != nil {
		return nil, err
	}

	parsed := resolveEnvPlaceholders(string(raw))

	cfg := &Config{}

	if err := yaml.Unmarshal([]byte(parsed), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
