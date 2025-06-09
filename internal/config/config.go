package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTl time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	DBUrl    string
	Secret   []byte
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("Файл config пуст")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Файл config ней найден: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Ошибка парсинга config файла: " + err.Error())
	}

	var secret []byte
	user, pass, host, dbName, secret := loadEnv()
	cfg.DBUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, dbName)
	cfg.Secret = secret

	return &cfg
}

func loadEnv() (user, pass, dbHost, dbName string, secret []byte) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	user = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbName = os.Getenv("DB_NAME")
	secret = []byte(os.Getenv("SECRET"))

	return
}

func fetchConfigPath() string {
	const op = "config.fetchConfigPath"

	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, fmt.Errorf("Не найден .env файл: %w", err)))
	}

	return os.Getenv("CONFIG_PATH")
}
