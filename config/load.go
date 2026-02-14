package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	return Config{
		Database: GetDatabaseEnv(),
		JWT:      GetJWTEnv(),
		Server:   GetServerEnv(),
		AiAgent:  GetAiAgentEnv(),
	}
}

func GetServerEnv() Server {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	port = ":" + port

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}

	return Server{
		Port:       port,
		CorsOrigin: corsOrigin,
	}
}
func GetJWTEnv() JWT {
	var jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET invalid")
	}
	return JWT{
		SecretKey: jwtSecret,
	}
}
func GetDatabaseEnv() Database {
	get := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("❌ %s invalid", key)
		}
		return val
	}

	return Database{
		Host:     get("DB_HOST"),
		User:     get("DB_USER"),
		Password: get("DB_PASSWORD"),
		DBName:   get("DB_NAME"),
		Port:     get("DB_PORT"),
		SSLMode:  get("DB_SSLMODE"),
		Schema:   get("DB_SCHEMA"),
	}
}

func GetAiAgentEnv() AiAgent {
	get := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("❌ %s invalid", key)
		}
		return val
	}

	model := os.Getenv("AI_AGENT_MODEL")
	if model == "" {
		model = "google/gemini-2.0-flash-exp:free"
	}

	return AiAgent{
		ApiKey: get("AI_AGENT_API_KEY"),
		Url:    get("AI_AGENT_URL"),
		Model:  model,
	}
}
