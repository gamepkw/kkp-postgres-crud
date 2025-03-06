package config

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Log     Log    `mapstructure:"log" validate:"required"`
	App     App    `mapstructure:"app" validate:"required"`
	Server  Server `mapstructure:"server" validate:"required"`
	Secrets Secret `mapstructure:"secrets" validate:"required"`
	DB      DB     `mapstructure:"postgres" validate:"required"`
}

type Log struct {
	Level string `mapstructure:"level" validate:"required"`
	Env   string `mapstructure:"env" validate:"required"`
}

type App struct {
	Name      string `mapstructure:"name" validate:"required"`
	ProjectID string `mapstructure:"project-id" validate:"required"`
}

type Server struct {
	Port     string `mapstructure:"port" validate:"required"`
	TimeZone string `mapstructure:"time-zone" validate:"required"`
}

type Secret struct {
	PostgresHost     string `envconfig:"SECRET_POSTGRES_HOST"`
	PostgresPort     int    `envconfig:"SECRET_POSTGRES_PORT"`
	PostgresUser     string `envconfig:"SECRET_POSTGRES_USER"`
	PostgresPassword string `envconfig:"SECRET_POSTGRES_PASSWORD"`
}

type DB struct {
	Name                string        `mapstructure:"db_name"`
	SSLMode             string        `mapstructure:"sslmode"`
	MaxOpenConns        *int          `mapstructure:"max_open_conns"`
	MaxIdleConns        *int          `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeHour time.Duration `mapstructure:"conn_max_lifetime_hour"`
}

type Masking []string

func Load(ctx context.Context) *AppConfig {
	var conf *AppConfig
	log.SetOutput(os.Stdout)
	configPath, ok := os.LookupEnv("API_CONFIG_PATH")
	if !ok {
		log.Println(ctx, "API_CONFIG_PATH not found, using default config")
		configPath = "./config"
	}
	configName, ok := os.LookupEnv("API_CONFIG_NAME")
	if !ok {
		log.Println(ctx, "API_CONFIG_NAME not found, using default config")
		configName = "config"
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(ctx, "config file not found. using default/env config: "+err.Error())
	}
	viper.AutomaticEnv()
	if err := viper.MergeConfig(strings.NewReader(viper.GetString("configs"))); err != nil {
		log.Panic(err.Error())
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Panic(err, "unmarshal config to struct error")
	}
	for _, value := range os.Environ() {
		pair := strings.SplitN(value, "=", 2)
		if strings.Contains(pair[0], "SECRET_") {
			keys := strings.Replace(pair[0], "SECRET_", "secrets.", -1)
			keys = strings.Replace(keys, "_", ".", -1)
			newKey := strings.Trim(keys, " ")
			newValue := strings.Trim(pair[1], " ")
			viper.Set(newKey, newValue)
		}
	}
	if err := godotenv.Load("./config/secret.env"); err != nil {
		log.Println("can't load ./config/secret.env", err)
	}
	if err := envconfig.Process("SECRET", &conf.Secrets); err != nil {
		log.Println("can't process SECRET", err)
	}
	return conf
}
