package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurtidev/kmf/internal/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRepository_Save(t *testing.T) {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbUser = "postgres"
	dbPassword = "postgres"
	dbName = "postgres"
	dbHost = "localhost"
	dbPort = "5432"

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	// Создание временной базы данных для тестов
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatal("Failed to connect to the test database:", err)
	}
	defer pool.Close()

	repo := &Repository{DB: pool}

	currency := &model.Currency{
		Title: "Test Currency",
		Code:  "TST",
		Value: 123.45,
		Date:  time.Now(),
	}

	err = repo.Save(context.Background(), currency)
	assert.NoError(t, err, "Error saving currency")
}

func TestRepository_FindByDateAndCode(t *testing.T) {
	// Создание временной базы данных для тестов
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbUser = "postgres"
	dbPassword = "postgres"
	dbName = "postgres"
	dbHost = "localhost"
	dbPort = "5432"

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	// Создание временной базы данных для тестов
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatal("Failed to connect to the test database:", err)
	}
	defer pool.Close()

	repo := &Repository{DB: pool}

	date := time.Now()
	code := "USD"

	currencies, err := repo.FindByDateAndCode(context.Background(), date, code)
	assert.NoError(t, err, "Error finding currencies by date and code")
	assert.NotNil(t, currencies, "Currencies not found")
}
