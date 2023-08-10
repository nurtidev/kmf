package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurtidev/kmf/internal/config"
	"github.com/nurtidev/kmf/internal/handler"
	"github.com/nurtidev/kmf/internal/repo/postgres"
	"github.com/nurtidev/kmf/internal/usecase"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %s", err)
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	// Создание пула подключений
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer pool.Close()

	repo := &postgres.Repository{DB: pool}
	useCase := &usecase.UseCase{
		Repo: repo,
	}
	currencyHandler := &handler.Handler{UseCase: useCase}

	r := mux.NewRouter()
	r.HandleFunc("/currency/save/{date}", currencyHandler.SaveCurrency).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", currencyHandler.GetCurrency).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Application is running...")
	log.Fatal(http.ListenAndServe(cfg.Port, r))

	// Теперь ваш конфиг загружен и валидирован, и вы можете продолжить с остальным кодом

}
