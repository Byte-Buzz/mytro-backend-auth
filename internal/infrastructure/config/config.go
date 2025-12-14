package config

import "time"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	CORS     CORSConfig
	Keys     TokenConfig
}

type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type LoggingConfig struct {
	Level    string
	Format   string
	Output   string
	FilePath string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type TokenConfig struct {
	PublicTokenLifetime  time.Duration
	PrivateTokenLifetime time.Duration

	privateKey string
	publicKey  string
}

func (k *TokenConfig) PrivateKey() string {
	key := k.privateKey
	k.privateKey = ""
	return key
}

func (k *TokenConfig) PublicKey() string {
	key := k.publicKey
	k.publicKey = ""
	return key
}

func LoadFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
