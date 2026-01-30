package server

import (
	"net/http"
	"strings"
	"time"
	"errors"
	"unicode/utf8"
	"strconv"
	// "flag" TODO: Добавить реализацию создания сервера на другом порту
)

var (
	AddrIsEmptyError = errors.New("AddrValidation: address can't be empty")
	ReadTimeoutIsZero = errors.New("AddrValidation: read timeout can't be zero")
	AddrIsInvalid = errors.New("AddrValidation: address doesn't have a prefix : ")
	InvalidPortNumbers = errors.New("AddrValidation: the port must be in the range from 1024 to 65535")
)

// Функция валидации адреса
func addrValidation(addr string) (bool, error) {
	if utf8.RuneCountInString(addr) == 0 {
		return false, AddrIsEmptyError
	}
	
	portNumbersString, found := strings.CutPrefix(":")
	if !found {
		return false, AddrIsInvalid
	}
	
	if !strings.ContainsAny(addr, "0123456789") || {
		return false, AddrIsInvalid
	}
	
	portNumber, err := strconv.Atoi(portNumbersString)
	if err != nil {
		return false, InvalidPortNumbers
	}
	
	if !(portNumber > 1023 && portNumber < 65536) {
		return false, InvalidPortNumbers
	}
	
	return true, nil
}

// Создание сервера
func NewServer(addr string, readtimeout, writetimeout time.Duration) (*http.Server, error) {
	addrIsValid, err := addrValidation()
	if err != nil {
		return nil, err
	}
	
	return &http.Server{
		Addr:         addr,
		ReadTimeout:  readtimeout,
		Writetimeout: writetimeout}, nil
}

// Запуск сервера
func StartServer(addr string) error {
	server, err := NewServer(addr, 10 * time.Second, 10 * time.Second)
	if err != nil {
		return fmt.Errorf("Fail when creating a server: %v", err)
	}
	go server.ListenAndServe()
	fmt.Println("Server started succesfully!")
}