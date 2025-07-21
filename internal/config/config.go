package config

import (
	"os"
	"time"
)

const (
	TheTVDBLoginURL  = "https://api4.thetvdb.com/v4/login"
	TheTVDBSearchURL = "https://api4.thetvdb.com/v4/search?query=%s&type=series"
	DefaultPort      = ":8080"
	JWTSecret        = "supersecretkey" // Change in production
	HTTPTimeout      = 10 * time.Second
)

var (
	PostgresURL = os.Getenv("POSTGRES_URL")
)
