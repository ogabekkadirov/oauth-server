package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// debug, info, warn, error,
	LogLevel         string
	HttpPort         string
	GrpcPort         string
	// db
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	// redis
	RedisAddr        string
	// sms
	SmsProvideApiKey string
	// jwt
	JWTPrivetKeyPath string
	JWTPublicKeyPath string
	JWTExpiresInSec  int
	// kafka
	KafkaAddress	string
	KafkaClientId	string
	KafkaGroupId	string
}

func Load() (Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	var config Config
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("HTTP_PORT", ":3030")
	v.SetDefault("GRPC_PORT", ":5000")
	v.SetDefault("POSTGRES_HOST", "localhost")
	v.SetDefault("POSTGRES_PORT", 5432)
	v.SetDefault("POSTGRES_USER", "postgres")
	v.SetDefault("POSTGRES_PASSWORD", "postgres")
	v.SetDefault("POSTGRES_DATABASE", "auth_db")
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("SMS_PROVIDER_API_KEY", "")
	v.SetDefault("JWT_PRIVETKEY_PATH", "")
	v.SetDefault("JWT_PUBLICKEY_PATH", "")
	v.SetDefault("JWT_EXPIRES_IN_SEC", 2630000) // 1 month
	v.SetDefault("KAFKA_ADDRESS", "broker:9092") 
	v.SetDefault("KAFKA_CLIENT_ID", "restaurant-restaurant-svc-consumer") 
	v.SetDefault("KAFKA_GROUP_ID", "restaurant-restaurant-svc-group") 

	config.LogLevel = v.GetString("LOG_LEVEL")
	config.HttpPort = v.GetString("HTTP_PORT")
	config.GrpcPort = v.GetString("GRPC_PORT")
	config.PostgresHost = v.GetString("POSTGRES_HOST")
	config.PostgresPort = v.GetInt("POSTGRES_PORT")
	config.PostgresUser = v.GetString("POSTGRES_USER")
	config.PostgresPassword = v.GetString("POSTGRES_PASSWORD")
	config.PostgresDatabase = v.GetString("POSTGRES_DATABASE")
	config.RedisAddr = v.GetString("REDIS_ADDR")
	config.SmsProvideApiKey = v.GetString("SMS_PROVIDER_API_KEY")
	config.JWTPrivetKeyPath = v.GetString("JWT_PRIVETKEY_PATH")
	config.JWTPublicKeyPath = v.GetString("JWT_PUBLICKEY_PATH")
	config.JWTExpiresInSec = v.GetInt("JWT_EXPIRES_IN_SEC")
	config.KafkaAddress = v.GetString("KAFKA_ADDRESS")
	config.KafkaClientId = v.GetString("KAFKA_CLIENT_ID")
	config.KafkaGroupId = v.GetString("KAFKA_GROUP_ID")

	return config, nil
}

func (c Config) NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	level, err := zapcore.ParseLevel(c.LogLevel)
	if err != nil {
		return nil, err
	}

	config.Level.SetLevel(level)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
