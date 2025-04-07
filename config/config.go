package config

import "os"

var (
	JWtSecret       string
	Sqlite_database string
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig() {
	JWtSecret = getEnv("JWT_SECRET", "make sure to provide JWT_SECRET")
	Sqlite_database = getEnv("SQLITE_DATABASE", "./database.db")
}
