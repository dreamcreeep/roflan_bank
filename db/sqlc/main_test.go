package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testStore Store
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Загружаем переменные окружения из .env файла в корне проекта
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file for tests")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// Формируем строку подключения для тестов
	// Обратите внимание: мы подключаемся к localhost, а не к 'postgres',
	// так как тесты обычно запускаются на хост-машине, а не в Docker.
	dbSource := fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPassword, dbName)

	testDB, err = sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
