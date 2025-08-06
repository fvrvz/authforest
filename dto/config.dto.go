package dto

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
	TimeZone string `yaml:"time_zone"`
}

type JWT struct {
	ExpiryMinutes int    `yaml:"expiry_minutes"`
	JWTSecret     string `yaml:"jwt_secret"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`
}
