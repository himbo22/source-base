package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server   ServerConfig     `yaml:"server"`
	Logger   LoggerConfig     `yaml:"logger"`
	Postgres PostgreSQLConfig `yaml:"postgres"`
	Redis    RedisConfig      `yaml:"redis"`
	Auth     AuthConfig       `yaml:"auth"`
	Email    EmailConfig      `yaml:"email"`
	Otel     TelemetryConfig  `yaml:"otel"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	OpenAPIPath     string `yaml:"openapiPath"`
	SwaggerPath     string `yaml:"swaggerPath"`
	ErrorStack      bool   `yaml:"errorStack"`
	ErrorLogEnabled bool   `yaml:"errorLogEnabled"`
	ErrorLogPattern string `yaml:"errorLogPattern"`
	LogLevel        string // Derived from logger config
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Path                 string   `yaml:"path"`
	File                 string   `yaml:"file"`
	Prefix               string   `yaml:"prefix"`
	Level                string   `yaml:"level"`
	TimeFormat           string   `yaml:"timeFormat"`
	CtxKeys              []string `yaml:"ctxKeys"`
	Header               bool     `yaml:"header"`
	StSkip               int      `yaml:"stSkip"`
	Stdout               bool     `yaml:"stdout"`
	RotateSize           int      `yaml:"rotateSize"`
	RotateExpire         int      `yaml:"rotateExpire"`
	RotateBackupLimit    int      `yaml:"rotateBackupLimit"`
	RotateBackupExpire   int      `yaml:"rotateBackupExpire"`
	RotateBackupCompress int      `yaml:"rotateBackupCompress"`
	RotateCheckInterval  string   `yaml:"rotateCheckInterval"`
	StdoutColorDisabled  bool     `yaml:"stdoutColorDisabled"`
	WriterColorEnable    bool     `yaml:"writerColorEnable"`
	Flags                int      `yaml:"flags"`
}

// PostgreSQLConfig holds database configuration
type PostgreSQLConfig struct {
	Host            string               `yaml:"host"`
	Port            int                  `yaml:"port"`
	User            string               `yaml:"user"`
	Password        string               `yaml:"password"`
	DBName          string               `yaml:"dbname"`
	SSLMode         string               `yaml:"sslmode"`
	MaxOpenConns    int                  `yaml:"max_open_conns"`
	MaxIdleConns    int                  `yaml:"max_idle_conns"`
	ConnMaxLifetime int                  `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int                  `yaml:"conn_max_idle_time"`
	Logger          DatabaseLoggerConfig `yaml:"logger"`
	Debug           bool                 `yaml:"debug"`
	Timezone        string               `yaml:"timezone"`
}

type DatabaseLoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Addrs           []string `yaml:"addrs"`
	MasterName      string   `yaml:"master_name"`
	Host            string   `yaml:"host"`
	Password        string   `yaml:"password"`
	Port            int      `yaml:"port"`
	Database        int      `yaml:"database"`
	PoolSize        int      `yaml:"pool_size"`
	MinIdleConns    int      `yaml:"min_idle_conns"`
	PoolTimeout     int      `yaml:"pool_timeout"`
	DialTimeout     int      `yaml:"dial_timeout"`
	ReadTimeout     int      `yaml:"read_timeout"`
	WriteTimeout    int      `yaml:"write_timeout"`
	MaxRetries      int      `yaml:"max_retries"`
	MaxRetryBackoff int      `yaml:"max_retry_backoff"`
	MinRetryBackoff int      `yaml:"min_retry_backoff"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JwtKeyId                 string `yaml:"jwtKeyId"` // Kong credential key (iss claim), e.g. bitzap-key
	AccessTokenExpireMinute  int    `yaml:"accessTokenExpireMinute"`
	RefreshTokenExpireMinute int    `yaml:"refreshTokenExpireMinute"`
	GracePeriodExpireSecond  int    `yaml:"gracePeriodExpireSecond"`
	GoogleClientId           string `yaml:"googleClientId"`
	RsaPublicKey             string `yaml:"rsaPublicKey"`
	RsaPrivateKey            string `yaml:"rsaPrivateKey"`
}

type EmailConfig struct {
	MailjetAPIKey    string `yaml:"mailjet_api_key" env:"MAILJET_API_KEY"`
	MailjetSecretKey string `yaml:"mailjet_secret_key" env:"MAILJET_SECRET_KEY"`
	FromEmail        string `yaml:"from_email" env:"FROM_EMAIL"`
	FromName         string `yaml:"from_name" env:"FROM_NAME"`
	AppURL           string `yaml:"app_url" env:"APP_URL"`
}

type TelemetryConfig struct {
	ServiceName    string  `mapstructure:"service_name" yaml:"service_name"`
	ServiceVersion string  `mapstructure:"service_version" yaml:"service_version"`
	Environment    string  `mapstructure:"environment" yaml:"environment"` // production | staging | development
	ExporterType   string  `mapstructure:"exporter_type" yaml:"exporter_type"`
	Endpoint       string  `mapstructure:"endpoint" yaml:"endpoint"`       // localhost:4317 or localhost:4318
	SampleRate     float64 `mapstructure:"sample_rate" yaml:"sample_rate"` // 1.0 = 100%, 0.1 = 10%
	Insecure       bool    `mapstructure:"insecure" yaml:"insecure"`       // disable TLS
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.Server.LogLevel = config.Logger.Level

	// Load sensitive data from environment variables

	return &config, nil
}
