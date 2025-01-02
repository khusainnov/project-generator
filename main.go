package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Получаем имя проекта из аргументов командной строки
	if len(os.Args) < 2 {
		fmt.Println("Укажите имя проекта: ./projectgen <projectname>")
		os.Exit(1)
	}
	projectName := os.Args[1]

	// Генерация структуры проекта
	err := generateProjectStructure(projectName)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Проект %s успешно создан!\n", projectName)
}

func generateProjectStructure(projectName string) error {
	// Структура директорий
	dirs := []string{
		"app/handler",
		"app/helpers",
		"app/config",
		"app/model",
		"app/processor",
		"app/repository",
		"app/infra/server",
		"app/infra/storage",
		"deployment",
		"specs/client",
		"specs/server/model",
	}

	// Шаблоны файлов
	files := map[string]string{
		"app/handler/handler.go": fmt.Sprintf(`package handler

import (
	"context"
	
	"go.uber.org/zap"

	"github.com/khusainnov/%s/app/processor"
)

// processor interface example
type Processor interface {
	Echo(ctx context.Context) error
}

// Handler example
type Handler struct {
	proc Processor
	log *zap.Logger
}

func New(log *zap.Logger, proc *processor.Processor) *Handler {
	return &Handler{
		log: log,
		proc: proc,
	}
}
`, projectName),
		"app/processor/processor.go": fmt.Sprintf(`package processor

import (
	"go.uber.org/zap"
)

// Processor example
type Processor struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Processor {
	return &Processor{
		log: log,
	}
}

func (p *Processor) Echo(ctx context.Context) error {
	p.log.Warn("implement  me")

	return nil
}

`),
		"app/config/config.go": `package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config example
type Config struct {
	Server *Server
	DB     *DB
}

type Server struct {
	Addr string ` + "`" + `envconfig:"ADDR"` + "`" + `
}

type DB struct {
	Host         string        ` + "`" + `envconfig:"HOST"` + "`" + `
	Port         string        ` + "`" + `envconfig:"PORT"` + "`" + `
	User         string        ` + "`" + `envconfig:"USER"` + "`" + `
	Password     string        ` + "`" + `envconfig:"PASSWORD"` + "`" + `
	Name         string        ` + "`" + `envconfig:"NAME"` + "`" + `
	SSLMode      string        ` + "`" + `envconfig:"SSL_MODE"` + "`" + `
	PingInterval time.Duration ` + "`" + `envconfig:"PING_INTERVAL"` + "`" + `
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess(Prefix, &c)

	return &c
}
`,
		"app/infra/storage/connect_db.go": fmt.Sprintf(`package infra

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/khusainnov/driver/postgres"
	"gitlab.com/khusainnov/%s/app/config"
	"go.uber.org/zap"
)

type ClientImpl struct {
	*sqlx.DB
}

func New(log *zap.Logger, cfg *config.DB) (*ClientImpl, error) {
	db, err := postgres.NewPostgresDB(
		postgres.ConfigPG{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			DBName:   cfg.Name,
			SSLMode:  cfg.SSLMode,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %w", err)
	}

	go func() {
		t := time.NewTicker(cfg.PingInterval)

		for range t.C {
			if err := db.Ping(); err != nil {
				log.Warn("failed to ping db", zap.Error(err))
			}
		}
	}()

	return &ClientImpl{db}, nil
}

func (db *ClientImpl) GetDB() *sqlx.DB {
	return db.DB
}

`, projectName),
		"app/infra/server/server.go": fmt.Sprintf(`package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"gitlab.com/khusainnov/%s/app/handler"
	"gitlab.com/khusainnov/%s/app/config"
	"go.uber.org/zap"
)

const (
	HTTP_API_PREFIX = "/jsonrpc/v2"
)

// Server example
type Server struct {
	httpServer *http.Server
	cfg        *config.Server
}

func New(cfg *config.Server) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Init(handlers *handler.Handler) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json2.NewCodec(), "application/json")
	if err := rpcServer.RegisterService(handlers, ""); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	router := mux.NewRouter()
	router.Handle(HTTP_API_PREFIX, rpcServer)

	s.httpServer = &http.Server{
		Addr:    s.cfg.Addr,
		Handler: router,
	}

	if err := http.ListenAndServe(s.cfg.Addr, router); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown the server: %w", err)
	}

	return nil
}

`, projectName, projectName),
		"app/app.go": fmt.Sprintf(`package app

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gitlab.com/khusainnov/%s/app/handler"
	"gitlab.com/khusainnov/%s/app/config"
	"gitlab.com/khusainnov/%s/app/infra/server"
	"gitlab.com/khusainnov/%s/app/infra/storage"
	"gitlab.com/khusainnov/%s/app/processor"
	"gitlab.com/khusainnov/%s/app/repository"
)

// App struct example
type App struct {
	Cfg    *config.Config
	Log    *zap.Logger
	Server *server.Server
	done   chan struct{}
}

func New(cfg *config.Config) *App {
	cfg := config.NewFromEnv()

	app := &App{
		Cfg: cfg,
	}

	log, err := app.CreateLogger()
	if err != nil {
		panic(err)
	}

	app.Log = log

	app.Log.Info("config init", zap.Any("config", *cfg))

	app.Log.Info("connecting to db")
	conn, err := storage.New(app.Log, cfg.DB)
	if err != nil {
		app.Log.Error("failed to connect to db")
		panic(err)
	}

	app.Log.Info("creating new handler")
	handlers := handler.New(
		app.Log,
		processor.New(app.Log),
	)

	app.Log.Info("creating new server")
	srv := server.New(app.Cfg.Server)
	app.Server = srv
	if err := app.Server.Init(handlers); err != nil {
		panic(err)
	}

	return app
}

func (a *App) CreateLogger() (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logCfg.EncoderConfig.TimeKey = "timestamp"

	return logCfg.Build(zap.AddCaller())
}

func (a *App) Run() {
	if a.Server != nil {
		go func() {
			defer func() { a.done <- struct{}{} }()
			a.Log.Info("starting server")
			err := a.Server.Run()
			if err != nil {
				a.Log.Error(err.Error(), zap.Error(err))
			}
		}()
	}

	a.Log.Info("%s running!")
	<-a.done
	a.stop()
}

func (a *App) stop() {
	a.Log.Info("%s stopping...")

	if err := a.Server.Stop(context.Background()); err != nil {
		a.Log.Error(err.Error(), zap.Error(err))
	}
}
`, projectName, projectName, projectName, projectName, projectName, projectName, projectName, projectName),
		"deployment/local.env": `# Local environment variables
APP_NAME=MyApp
`,
		"deployment/docker-compose.yaml": `version: "3"
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - APP_NAME=MyApp
    ports:
      - "8080:8080"
`,
		"main.go": fmt.Sprintf(`package main

import (
	"gitlab.com/khusainnov/%s/app"
)

func main() {
	inviteapp := app.New()
	inviteapp.Run()
}
`, projectName),
		"go.mod": "module " + projectName + "\n\ngo 1.23",
		".gitignore": `# Binaries
*.exe
*.out
.idea

# Logs
*.log

# OS generated
.DS_Store
`,
		"Dockerfile": `# Stage 1: Build
FROM golang:1.23 AS builder

# Making working dir
WORKDIR /app

# Copy go.mod and go.sum and installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Building binary
RUN go build -o main .

# Stage 2: Run
FROM debian:bullseye-slim

# Installing dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Making working dir
WORKDIR /app

# Copy binary file from first stage
COPY --from=builder /app/main .

# App runner
CMD ["./main"]
`,
		"README.md": "# " + projectName + "\n\nProject description.",
	}

	// Creating dirs
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectName, dir), 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}

	// Creating files
	for path, content := range files {
		filePath := filepath.Join(projectName, path)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	return nil
}
