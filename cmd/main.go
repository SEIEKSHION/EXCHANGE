package main

import (
	"fmt"

	"github.com/SEIEKSHION/Exchanger/internal/handlers"
	"github.com/SEIEKSHION/Exchanger/internal/models"
	"github.com/SEIEKSHION/Exchanger/internal/server"
)

func main() {
	body, err := models.GetVaultExchange()
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}

	valutes, err := models.ProceedExchangeVaults(body)
	if err != nil {
		panic(fmt.Errorf("Main:\n\t%v", err))
	}

	err = models.PrintValutes(valutes)
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}

	handler := handlers.NewHandler(valutes)

	// Создание и запуск сервера
	srv, err := server.NewServer(":8080", handler)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		return
	}

	if err := srv.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	// Здесь можно выполнить другие задачи
	// или использовать сигналы для graceful shutdown

	// Пример ожидания сигнала завершения
	select {}
}
