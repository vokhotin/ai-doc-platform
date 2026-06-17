package config

import "os"

type Config struct {
	Port         string
	UploadDir    string
	DatabaseURL  string
	InferenceURL string
}

func Load() Config {
	return Config{
		Port:         getEnv("GATEWAY_SERVICE_PORT", "8080"),
		UploadDir:    getEnv("GATEWAY_SERVICE_UPLOAD_DIR", "./uploads"),
		DatabaseURL:  getEnv("GATEWAY_SERVICE_DB_URL", ""),
		InferenceURL: getEnv("INFERENCE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
