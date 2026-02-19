package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/SEIEKSHION/Exchanger/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	AddrIsEmptyError   = errors.New("addr validation: address can't be empty")
	ReadTimeoutIsZero  = errors.New("addr validation: read timeout can't be zero")
	AddrIsInvalid      = errors.New("addr validation: address doesn't have a prefix : ")
	InvalidPortNumbers = errors.New("addr validation: the port must be in the range from 1024 to 65535")
)

type Server struct {
	httpServer *http.Server
	mu         sync.Mutex
	running    bool
}

// Валидация адреса
func addrValidation(addr string) error {
	if utf8.RuneCountInString(addr) == 0 {
		return AddrIsEmptyError
	}

	portNumbersString, found := strings.CutPrefix(addr, ":")
	if !found {
		return AddrIsInvalid
	}

	if !strings.ContainsAny(addr, "0123456789") {
		return AddrIsInvalid
	}

	portNumber, err := strconv.Atoi(portNumbersString)
	if err != nil {
		return InvalidPortNumbers
	}

	if !(portNumber > 1023 && portNumber < 65536) {
		return InvalidPortNumbers
	}

	return nil
}

// Создание сервера
func NewServer(addr string, measurementhandler *handlers.MeasurementHandler, exchangehandler *handlers.Handler) (*Server, error) {
	if err := addrValidation(addr); err != nil {
		return nil, fmt.Errorf("server creation failed: %w", err)
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Статические файлы
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/favicon.png", "./static/images/favicon.png")

	// Маршруты
	api := r.Group("/api")
	{
		api.GET("/valutes", exchangehandler.GetValutes)
		api.POST("/convert", exchangehandler.ConvertCurrency)
		api.GET("/createmesaurement", measurementhandler.Create)
	}

	r.GET("/", exchangehandler.MainPage)
	r.GET("/exchanger", exchangehandler.Exchanger)
	// r.GET("/muscles", musclehandler.MusclesPage)  покрыть миддлварем с авторизацией
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: srv,
	}, nil
}

// Запуск сервера (неблокирующий)
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("server is already running")
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server error: %v\n", err)
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	fmt.Printf("Server started successfully on %s!\n", s.httpServer.Addr)
	return nil
}

// Остановка сервера
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return errors.New("server is not running")
	}
	s.running = false
	s.mu.Unlock()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("Server stopped successfully!")
	return nil
}

// Проверка состояния сервера
func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Запуск сервера (старый интерфейс для совместимости)
func StartServer(addr string, exchangehandler *handlers.Handler, muscleshandler *handlers.MusclesHandler) error {
	server, err := NewServer(addr, muscleshandler, exchangehandler)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return server.Start()
}
