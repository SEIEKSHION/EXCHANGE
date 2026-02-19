package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SEIEKSHION/Exchanger/internal/handlers"
	"github.com/SEIEKSHION/Exchanger/internal/models"
	"github.com/SEIEKSHION/Exchanger/internal/repository"
	"github.com/SEIEKSHION/Exchanger/internal/server"
)

func main() {

	// получение значений флагов
	portPtr := flag.Int("port", 8080, "The server will be started on this port")
	flag.Parse()

	// получение курса валют
	body, err := models.GetVaultExchange()
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}

	// обработка  курсов
	valutes, err := models.ProceedExchangeVaults(body)
	if err != nil {
		panic(fmt.Errorf("Main:\n\t%v", err))
	}

	connString := "host=localhost port=5432 user=seiekshion password=UnderMind35327711_ dbname=exchanger sslmode=disable"

	// Инициализация БД
	dbClient, err := repository.NewClient(connString)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer dbClient.DB.Close()

	// инициализация хэндлеров
	// muscleHandler := handlers.NewMuscleHandler()

	measurementRepo := repository.NewMeasurementRepository(dbClient)
	measurementHandler := handlers.NewMeasurementHandler(measurementRepo)
	exchangeHandler := handlers.NewHandler(valutes)

	// Создание и запуск сервера
	srv, err := server.NewServer(fmt.Sprintf(":%d", *portPtr), measurementhandler, exchangeHandler) // передаём порт
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// обработка при запуске сервера
	if err := srv.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

	// Настройка graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// обработка при остановке сервера
	if err := srv.Stop(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
		os.Exit(1)
	}

}
