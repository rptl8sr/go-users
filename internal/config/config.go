package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

const (
	defaultHTTPPort  = 8080
	defaultSwaggerUI = "/swagger"
)

// App represents application-level configuration, including mode, debug state, and Swagger UI settings.
type App struct {
	Mode      string `env:"APP_MODE" env-default:"prod"`
	Debug     bool   `env:"APP_DEBUG" envDefault:"false"`
	SwaggerUI string `env:"APP_SWAGGER_UI" envDefault:"/swagger"`
}

// HTTP represents the configuration for the HTTP server including port, host, and timeout settings.
type HTTP struct {
	Port         int    `env:"HTTP_PORT" env-default:"8080"`
	Host         string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	ReadTimeout  int    `env:"HTTP_READ_TIMEOUT" env-default:"5"`
	WriteTimeout int    `env:"HTTP_WRITE_TIMEOUT" env-default:"10"`
	IdleTimeout  int    `env:"HTTP_IDLE_TIMEOUT" env-default:"120"`
}

// Log represents logging configuration including log level, format, and output destination.
type Log struct {
	Level  slog.Level `env:"LOG_LEVEL" env-default:"info"`
	Format string     `env:"LOG_FORMAT" env-default:"json"`
	Output string     `env:"LOG_OUTPUT"`
}

// OpenAPI represents configuration related to OpenAPI specifications and routes.
type OpenAPI struct {
	SpecPath  string `env:"OPENAPI_SPEC_PATH" env-default:"openapi/openapi.yaml"`
	APIPrefix string `env:"OPENAPI_API_PREFIX" env-default:"/api/v1"`
}

// Database represents the configuration for a database connection, including host, port, credentials, and settings.
type Database struct {
	Host           string `env:"DB_HOST" env-default:"localhost"`
	Port           int    `env:"DB_PORT" env-default:"5432"`
	User           string `env:"DB_USER" env-default:"postgres"`
	Password       string `env:"DB_PASSWORD" env-default:"postgres"`
	Name           string `env:"DB_NAME" env-default:"users"`
	SSLMode        string `env:"DB_SSL_MODE" env-default:"disable"`
	MaxConnections int    `env:"DB_MAX_CONNECTIONS" env-default:"10"`
}

// Config represents the configuration structure for the application, including settings for App, HTTP, Log, OpenAPI, and Database.
type Config struct {
	App      App
	HTTP     HTTP
	Log      Log
	OpenAPI  OpenAPI
	Database Database
}

// New initializes a new Config object by reading environment variables and applying default settings and flags.
func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	setMode(cfg)

	return cfg, nil
}

// setMode configures the application settings based on the specified mode (dev, test, prod) in the provided Config object.
func setMode(cfg *Config) {
	switch cfg.App.Mode {
	case "dev":
		cfg.App.Debug = true
		cfg.Log.Level = slog.LevelDebug
		cfg.Log.Format = "json"
		cfg.App.SwaggerUI = defaultSwaggerUI
		cfg.Database.SSLMode = "disable"
		cfg.HTTP.Port = defaultHTTPPort
	case "test":
		cfg.App.Debug = true
		cfg.Log.Level = slog.LevelDebug
		cfg.Log.Format = "text"
		cfg.App.SwaggerUI = ""
		cfg.Database.SSLMode = "disable"
		cfg.HTTP.Port = defaultHTTPPort
	case "prod":
		cfg.App.Debug = false
		cfg.Log.Level = slog.LevelInfo
		cfg.Log.Format = "json"
		cfg.App.SwaggerUI = ""
		cfg.Database.SSLMode = "require"
		cfg.HTTP.Port = defaultHTTPPort
	}
}

// String converts the Config object into a formatted string representation summarizing HTTP, log, OpenAPI, and database settings.
func (c *Config) String() string {
	return fmt.Sprintf(
		"HTTP: %s:%d, Log: %s (%s), OpenAPI: %s, DB: %s:%d",
		c.HTTP.Host,
		c.HTTP.Port,
		c.Log.Level,
		c.Log.Format,
		c.OpenAPI.SpecPath,
		c.Database.Host,
		c.Database.Port,
	)
}
