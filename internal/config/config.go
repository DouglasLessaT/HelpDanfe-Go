package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config representa todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SEFAZ    SEFAZConfig
	Bank     BankConfig
	Log      LogConfig
	Cache    CacheConfig
}

// ServerConfig representa as configurações do servidor
type ServerConfig struct {
	Port      string
	Host      string
	Environment string
}

// DatabaseConfig representa as configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// SEFAZConfig representa as configurações da SEFAZ
type SEFAZConfig struct {
	Ambiente     string
	UF           string
	Timeout      time.Duration
	CertPath     string
	CertPassword string
}

// BankConfig representa as configurações bancárias
type BankConfig struct {
	Itau        BankAPIConfig
	Bradesco    BankAPIConfig
	OpenBanking BankAPIConfig
}

// BankAPIConfig representa as configurações de uma API bancária
type BankAPIConfig struct {
	URL         string
	ClientID    string
	ClientSecret string
	Timeout     time.Duration
}

// LogConfig representa as configurações de log
type LogConfig struct {
	Level string
	File  string
}

// CacheConfig representa as configurações de cache
type CacheConfig struct {
	TTL      time.Duration
	RedisURL string
}

// Load carrega as configurações das variáveis de ambiente
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "helpdanfe"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		SEFAZ: SEFAZConfig{
			Ambiente:     getEnv("SEFAZ_AMBIENTE", "homologacao"),
			UF:           getEnv("SEFAZ_UF", "SP"),
			Timeout:      getEnvDuration("SEFAZ_TIMEOUT", 30*time.Second),
			CertPath:     getEnv("CERT_PATH", "./certs/certificado.p12"),
			CertPassword: getEnv("CERT_PASSWORD", ""),
		},
		Bank: BankConfig{
			Itau: BankAPIConfig{
				URL:          getEnv("ITAÚ_API_URL", ""),
				ClientID:     getEnv("ITAÚ_CLIENT_ID", ""),
				ClientSecret: getEnv("ITAÚ_CLIENT_SECRET", ""),
				Timeout:      getEnvDuration("ITAÚ_TIMEOUT", 30*time.Second),
			},
			Bradesco: BankAPIConfig{
				URL:          getEnv("BRADESCO_API_URL", ""),
				ClientID:     getEnv("BRADESCO_CLIENT_ID", ""),
				ClientSecret: getEnv("BRADESCO_CLIENT_SECRET", ""),
				Timeout:      getEnvDuration("BRADESCO_TIMEOUT", 30*time.Second),
			},
			OpenBanking: BankAPIConfig{
				URL:          getEnv("OPEN_BANKING_URL", ""),
				ClientID:     getEnv("OPEN_BANKING_CLIENT_ID", ""),
				ClientSecret: getEnv("OPEN_BANKING_CLIENT_SECRET", ""),
				Timeout:      getEnvDuration("OPEN_BANKING_TIMEOUT", 30*time.Second),
			},
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "./logs/app.log"),
		},
		Cache: CacheConfig{
			TTL:      getEnvDuration("CACHE_TTL", time.Hour),
			RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
	}

	return config, nil
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvDuration obtém uma variável de ambiente como duração ou retorna um valor padrão
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvInt obtém uma variável de ambiente como inteiro ou retorna um valor padrão
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDatabaseDSN retorna a string de conexão do banco de dados
func (c *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

// IsProduction retorna true se o ambiente for de produção
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment retorna true se o ambiente for de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}
