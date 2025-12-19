package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port        string
	DBPath      string
	DataDir     string
	Timezone    string
	TemplateDir string
}

var (
	// CommitHash is the git commit hash at build time
	CommitHash = "unknown"
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./data/web2rss.sqlite3"),
		DataDir:     getEnv("DATA_DIR", "./data"),
		Timezone:    getEnv("APP_TIMEZONE", "UTC"),
		TemplateDir: getEnv("TEMPLATE_DIR", "templates"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
