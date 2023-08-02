package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/prinick96/elog"
)

// Default PORT to set if .env PORT var is empty or can't load
const DEFAULT_PORT_IF_EMPTY = "8080"

// Env config struct
type EnvApp struct {
	// Server Envs
	PORT string

	JWT_SECRET string
	// Database Envs
	DB_HOST       string
	DB_PORT       string
	DB_USERNAME   string
	DB_PASSWORD   string
	DB_SERVICE    string
	CIDI_PASS     string
	CIDI_KEY      string
	ID_APP        string
	BASE_CIDI_URI string
	GIN_MODE	  string
}

// Get the env configuration
func GetEnv(env_file string) EnvApp {
	err := godotenv.Load(env_file)
	elog.New(elog.PANIC, "Error loading "+env_file+" file", err)

	// Heroku smell
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT_IF_EMPTY
	}

	return EnvApp{
		PORT:          port,
		DB_HOST:       os.Getenv("DB_HOST"),
		DB_PORT:       os.Getenv("DB_PORT"),
		DB_USERNAME:   os.Getenv("DB_USERNAME"),
		DB_PASSWORD:   os.Getenv("DB_PASSWORD"),
		DB_SERVICE:    os.Getenv("DB_SERVICE"),
		JWT_SECRET:    os.Getenv("JWT_SECRET"),
		CIDI_PASS:     os.Getenv("CIDI_PASS"),
		CIDI_KEY:      os.Getenv("CIDI_KEY"),
		ID_APP:        os.Getenv("ID_APP"),
		BASE_CIDI_URI: os.Getenv("BASE_CIDI_URI"),
		GIN_MODE:      os.Getenv("GIN_MODE"),
	}
}
