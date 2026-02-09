package config

type Config struct {
	Database Database
	JWT      JWT
	Server   Server
	AiAgent  AiAgent
}

type Database struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	Schema   string
}

type JWT struct {
	SecretKey string
}

type Server struct {
	Port       string
	CorsOrigin string
}

type AiAgent struct {
	ApiKey string
	Url    string
}
