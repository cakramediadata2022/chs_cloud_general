package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	MySQL    MySQL
	Redis    RedisConfig
	CM       CM
	MongoDB  MongoDB
	Cookie   Cookie
	Store    Store
	Session  Session
	Metrics  Metrics
	Logger   Logger
	AWS      AWS
	Jaeger   Jaeger
	Discord  Discord
	Tada     Tada
}

// Server config struct
type ServerConfig struct {
	AppName           string
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	Staging           bool
	BaseURL           string
	Domain            string
	PublicPath        string
	SecretPath        string
	MultiDatabase     bool
}

// Discord Config
type Discord struct {
	Api string
	Run bool
}

// Discord Config
type Tada struct {
	Host        string
	ServiceName string
}

// Channel Manager Config
type CM struct {
	RGURL    string //RateGain
	CXURL    string //Channex
	Username string
	Password string
}

// Logger config
type Logger struct {
	LogFileEnabled    bool
	LogMaxSize        int
	LogMaxBackups     int
	LogMaxAge         int
	LogFilename       string
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// MongoDB config
type MongoDB struct {
	MongoURI string
}

type MySQL struct {
	MySqlHost     string
	MySqlPort     int
	MySqlUser     string
	MySqlPassword string
	MySqlDatabase string
	MaxOpenConns  int
	MaxIdleConns  int
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Store config
type Store struct {
	ImagesFolder string
}

// AWS S3
type AWS struct {
	Endpoint       string
	MinioAccessKey string
	MinioSecretKey string
	UseSSL         bool
	MinioEndpoint  string
}

// AWS S3
type Jaeger struct {
	Host         string
	ServiceName  string
	LogSpans     bool
	LogQuery     bool
	SamplerRatio float64
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
