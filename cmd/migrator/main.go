package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	user, pass, host, dbName, path := loadEnv()

	m, err := migrate.New(
		"file://"+path,
		fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, dbName),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Нет новых миграций")
			return
		}

		panic(err)
	}

	fmt.Println("Миграция выполнена")
}

func loadEnv() (user, pass, dbHost, dbName, migrationsPath string) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	user = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbName = os.Getenv("DB_NAME")
	migrationsPath = os.Getenv("MIGRATIONS_PATH")

	return
}
