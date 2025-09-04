package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseDriver string
	DatabaseAddr   string
	DatabaseUser   string
	DatabasePass   string
	DatabaseName   string
	RedisAddr      string
	PhotoQueryPath string
	MusicQueryPath string
	VideoQueryPath string
	PhotoPath      string
	MusicPath      string
	VideoPath      string
	ChunkPath      string
	FilePath       string
}

// AppConfig is the global configuration object
var AppConfig *Config

// getEnv reads an environment variable or returns a fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// init loads configuration from environment variables
func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	AppConfig = &Config{
		DatabaseDriver: getEnv("DATABASE_DRIVER", "mysql"),
		DatabaseAddr:   getEnv("DATABASE_ADDR", "127.0.0.1:3306"),
		DatabaseUser:   getEnv("DATABASE_USER", "root"),
		DatabasePass:   getEnv("DATABASE_PASS", "menchowgan"), // IMPORTANT: No default for password
		DatabaseName:   getEnv("DATABASE_NAME", "gmc_blog"),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		PhotoQueryPath: getEnv("PHOTO_QUERY_PATH", "http://localhost/"),
		MusicQueryPath: getEnv("MUSIC_QUERY_PATH", "http://localhost/"),
		VideoQueryPath: getEnv("VIDEO_QUERY_PATH", "http://localhost/"),
		PhotoPath:      getEnv("PHOTO_PATH", "/usr/local/share/images/"),
		MusicPath:      getEnv("MUSIC_PATH", "/usr/local/share/music/"),
		VideoPath:      getEnv("VIDEO_PATH", "/usr/local/share/video/"),
		ChunkPath:      getEnv("CHUNK_PATH", "/usr/local/share/chunks/"),
		FilePath:       getEnv("FILE_PATH", "/usr/local/share/files/"),
	}
}
