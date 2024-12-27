package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var (
	errConfigEnvEmpty    = errors.New("config path env is not set")
	errConfigFileInvalid = errors.New("config file has invalid format")
)

type EnvType string

const (
	EnvLocal   EnvType = "local"
	EnvDev     EnvType = "dev"
	EnvTesting EnvType = "testing"

	EnvConfigPath string = "APP_CONFIG_PATH"
)

type Config struct {
	Env             EnvType    `yaml:"env" env-required:"true"`
	ApplicationName string     `yaml:"application_name" env-required:"true"`
	GrpcServer      GrpcServer `yaml:"grpc_server" env-required:"true"`
	Services        Services   `yaml:"services" env-required:"true"`
	Timeouts        Timeouts   `yaml:"timeouts" env-required:"true"`
	Tracing         Tracing    `yaml:"tracing" env-required:"true"`
	Redis           Redis      `yaml:"redis" env-required:"true"`
}

type GrpcServer struct {
	Port    string `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout string `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT" env-required:"true"`
}

type Services struct {
	Authentication AuthenticationService `yaml:"authentication_service" env-required:"true"`
	Form           FormService           `yaml:"form_service" env-required:"true"`
	Submission     SubmissionService     `yaml:"submission_service" env-required:"true"`
}

type AuthenticationService struct {
	Host string `yaml:"host" env:"AUTHENTICATION_SERVICE_HOST" env-required:"true"`
	Port string `yaml:"port" env:"AUTHENTICATION_SERVICE_PORT" env-required:"true"`
}

type FormService struct {
	Host string `yaml:"host" env:"FORM_SERVICE_HOST" env-required:"true"`
	Port string `yaml:"port" env:"FORM_SERVICE_PORT" env-required:"true"`
}

type SubmissionService struct {
	Host string `yaml:"host" env:"SUBMISSION_SERVICE_HOST" env-required:"true"`
	Port string `yaml:"port" env:"SUBMISSION_SERVICE_PORT" env-required:"true"`
}

type Timeouts struct {
	EndpointExecutionTimeoutMS string `yaml:"endpoint_execution_timeout_ms" env:"ENDPOINT_EXECUTION_TIMEOUT_MS" env-required:"true"`
}

type Tracing struct {
	Host                string `yaml:"host" env:"JAEGER_HOST" env-required:"true"`
	Port                string `yaml:"port" env:"JAEGER_PORT" env-required:"true"`
	StatisticTracerName string `yaml:"statistic_tracer_name" env:"STATISTIC_TRACER_NAME" env-required:"true"`
}

type Redis struct {
	User           string `yaml:"user" env:"REDIS_USER" env-required:"true"`
	Password       string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`
	Host           string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port           string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
	DatabaseNumber string `yaml:"database_number" env:"REDIS_DATABASE_NUMBER" env-required:"true"`
	EntityTTL      string `yaml:"entity_ttl" env:"REDIS_ENTITY_TTL" env-required:"true"`
}

func Load() (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(Load).Pointer()).Name()

	pathToConfig := os.Getenv(EnvConfigPath)
	if pathToConfig == "" {
		return nil, fmt.Errorf("%s -> %w", op, errConfigEnvEmpty)
	}

	return loadByPath(pathToConfig)
}

func loadByPath(pathToConfig string) (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(loadByPath).Pointer()).Name()

	var config Config
	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s -> %w; %w", op, errConfigFileInvalid, err)
	}

	return &config, nil
}

func (cfg *Config) GetEndpointExecutionTimeout() time.Duration {
	op := runtime.FuncForPC(reflect.ValueOf(cfg.GetEndpointExecutionTimeout).Pointer()).Name()

	milliseconds, err := strconv.Atoi(cfg.Timeouts.EndpointExecutionTimeoutMS)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s, op=%s, err=%s", "unable to convert endpoint timeout duration", op, err.Error()))
	}

	return time.Duration(milliseconds) * time.Millisecond
}

func (cfg *Config) GetRedisEntityTTL() time.Duration {
	op := runtime.FuncForPC(reflect.ValueOf(cfg.GetRedisEntityTTL).Pointer()).Name()

	parsedTTL, err := time.ParseDuration(cfg.Redis.EntityTTL)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s, op=%s, err=%s", "unable to parse entity ttl into time.Duration", op, err.Error()))
	}

	return parsedTTL
}
