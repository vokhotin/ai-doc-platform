package config

import "os"

type Config struct {
	Port      string
	UploadDir string
}

func Load() Config {
	return Config{
		Port:      getEnv("GATEWAY_SERVICE_PORT", "8080"),
		UploadDir: getEnv("GATEWAY_SERVICE_UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
