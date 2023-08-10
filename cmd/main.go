package main

import (
	"flag"
	"fmt"
	"github.com/nurtidev/kmf/internal/config"
	"log"
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

	// Теперь ваш конфиг загружен и валидирован, и вы можете продолжить с остальным кодом
	fmt.Println("Application is running...")
}
