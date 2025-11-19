package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database   DatabaseConfig
	JWT        JWTConfig
	OpenRouter OpenRouterConfig
	Supabase   SupabaseConfig
	Server     ServerConfig
	CORS       CORSConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig holds JWT authentication settings
type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
	RefreshTime    time.Duration
}

// OpenRouterConfig holds OpenRouter API settings
type OpenRouterConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

// SupabaseConfig holds Supabase settings
type SupabaseConfig struct {
	URL       string
	AnonKey   string
	JWTSecret string
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

// CORSConfig holds CORS settings
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// LoadConfig loads configuration from environment variables and config files
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	setDefaults()

	// Enable environment variable override
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file (optional - will use env vars if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; using environment variables
	}

	var config Config

	// Database Config - check for DATABASE_URL first (Railway/Heroku style)
	databaseURL := viper.GetString("DATABASE_URL")
	if databaseURL != "" {
		// Parse DATABASE_URL
		dbConfig, err := parseDatabaseURL(databaseURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing DATABASE_URL: %w", err)
		}
		config.Database = *dbConfig
	} else {
		// Use individual environment variables
		config.Database = DatabaseConfig{
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			User:            viper.GetString("database.user"),
			Password:        viper.GetString("database.password"),
			DBName:          viper.GetString("database.dbname"),
			SSLMode:         viper.GetString("database.sslmode"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
		}
	}

	// JWT Config
	config.JWT = JWTConfig{
		Secret:         viper.GetString("jwt.secret"),
		ExpirationTime: viper.GetDuration("jwt.expiration_time"),
		RefreshTime:    viper.GetDuration("jwt.refresh_time"),
	}

	// OpenRouter Config
	config.OpenRouter = OpenRouterConfig{
		APIKey:  viper.GetString("openrouter.api_key"),
		BaseURL: viper.GetString("openrouter.base_url"),
		Model:   viper.GetString("openrouter.model"),
		Timeout: viper.GetDuration("openrouter.timeout"),
	}

	// Supabase Config
	config.Supabase = SupabaseConfig{
		URL:       viper.GetString("supabase.url"),
		AnonKey:   viper.GetString("supabase.anon_key"),
		JWTSecret: viper.GetString("supabase.jwt_secret"),
	}

	// Server Config
	config.Server = ServerConfig{
		Port:            viper.GetInt("server.port"),
		Host:            viper.GetString("server.host"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
		Environment:     viper.GetString("server.environment"),
	}

	// CORS Config
	config.CORS = CORSConfig{
		AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
		AllowedMethods:   viper.GetStringSlice("cors.allowed_methods"),
		AllowedHeaders:   viper.GetStringSlice("cors.allowed_headers"),
		ExposedHeaders:   viper.GetStringSlice("cors.exposed_headers"),
		AllowCredentials: viper.GetBool("cors.allow_credentials"),
		MaxAge:           viper.GetInt("cors.max_age"),
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	// JWT defaults
	viper.SetDefault("jwt.expiration_time", 24*time.Hour)
	viper.SetDefault("jwt.refresh_time", 7*24*time.Hour)

	// OpenRouter defaults
	viper.SetDefault("openrouter.base_url", "https://openrouter.ai/api/v1")
	viper.SetDefault("openrouter.model", "openai/gpt-4-turbo-preview")
	viper.SetDefault("openrouter.timeout", 30*time.Second)

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 15*time.Second)
	viper.SetDefault("server.write_timeout", 15*time.Second)
	viper.SetDefault("server.shutdown_timeout", 10*time.Second)
	viper.SetDefault("server.environment", "development")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Content-Type", "Authorization"})
	viper.SetDefault("cors.exposed_headers", []string{"Content-Length"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 3600)
}

// parseDatabaseURL parses a DATABASE_URL connection string
// Format: postgresql://user:password@host:port/dbname?sslmode=require
func parseDatabaseURL(dbURL string) (*DatabaseConfig, error) {
	// Remove postgresql:// or postgres:// prefix
	dbURL = strings.TrimPrefix(dbURL, "postgresql://")
	dbURL = strings.TrimPrefix(dbURL, "postgres://")

	// Split by @
	parts := strings.Split(dbURL, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid DATABASE_URL format")
	}

	// Parse credentials
	creds := strings.Split(parts[0], ":")
	if len(creds) != 2 {
		return nil, fmt.Errorf("invalid credentials in DATABASE_URL")
	}
	user := creds[0]
	password := creds[1]

	// Parse host/port/db
	hostParts := strings.Split(parts[1], "/")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("invalid host/database in DATABASE_URL")
	}

	// Parse host and port
	hostPort := strings.Split(hostParts[0], ":")
	host := hostPort[0]
	port := 5432 // default
	if len(hostPort) == 2 {
		var err error
		_, err = fmt.Sscanf(hostPort[1], "%d", &port)
		if err != nil {
			return nil, fmt.Errorf("invalid port in DATABASE_URL: %w", err)
		}
	}

	// Parse dbname and query params
	dbParts := strings.Split(hostParts[1], "?")
	dbname := dbParts[0]

	// Default to require for production
	sslmode := "require"
	if len(dbParts) == 2 {
		// Parse query params
		params := strings.Split(dbParts[1], "&")
		for _, param := range params {
			kv := strings.Split(param, "=")
			if len(kv) == 2 && kv[0] == "sslmode" {
				sslmode = kv[1]
			}
		}
	}

	return &DatabaseConfig{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		DBName:          dbname,
		SSLMode:         sslmode,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}, nil
}

// validateConfig validates required configuration values
func validateConfig(config *Config) error {
	// Validate Database
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate JWT
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	if len(config.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// OpenRouter and Supabase are optional - only validate if provided
	// This allows for basic deployment without these services

	// Validate Server
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}
