package config

import "os"

type Config struct {
	HTTPAddr       string
	DatabaseURL    string
	AIInferenceURL string
}

func Load() Config {
	return Config{
		HTTPAddr:       env("BACKEND_HTTP_ADDR", ":8080"),
		DatabaseURL:    env("DATABASE_URL", ""),
		AIInferenceURL: env("AI_INFERENCE_URL", "http://localhost:8000"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
