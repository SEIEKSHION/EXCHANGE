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
	"github.com/SEIEKSHION/Exchanger/internal/server"
)

func main() {

	// получение значений флагов
	portPtr := flag.Int("port", 8080, "The server will be started on this port")
	flag.Parse()
	fmt.Println(*portPtr)
	body, err := models.GetVaultExchange()
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}

	valutes, err := models.ProceedExchangeVaults(body)
	if err != nil {
		panic(fmt.Errorf("Main:\n\t%v", err))
	}
	muscle, err := models.NewMuscle("Бицепс", 94.6, time.Now())
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
	var muscles []models.Muscle = []models.Muscle{muscle}

	muscleHandler := handlers.NewMuscleHandler(muscles)
	exchangeHandler := handlers.NewHandler(valutes)

	// Создание и запуск сервера
	srv, err := server.NewServer(fmt.Sprintf(":%d", *portPtr), muscleHandler, exchangeHandler) // передаём порт
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

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

	if err := srv.Stop(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
		os.Exit(1)
	}

}
