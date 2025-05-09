package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	fmt.Printf("Project %s successfully created!!!\n", projectName)
}

func generateProjectStructure(projectName string) error {
	// Структура директорий
	dirs := []string{
		"app/handler",
		"app/helpers",
		"app/config",
		"app/model",
		"app/processor",
		"app/service",
		"app/repository",
		"app/infra/server",
		"app/infra/storage",
		"app/infra/log",
		"app/infra/system",
		"deployment",
		"specs/client",
		"specs/server/params/model",
		"specs/server/response/model",
	}

	// Шаблоны файлов
	files := map[string]string{
		"app/handler/echo.go": fmt.Sprintf(`package handler

import (
	"net/http"

	"go.uber.org/zap"

	reqParams "github.com/khusainnov/%s/specs/server/params/model"
	respParams "github.com/khusainnov/%s/specs/server/response/model"
)

func (h *Handler) Echo(r *http.Request, params *reqParams.EchoParams, resp *respParams.EchoResponse) error {
	ctx := r.Context()
	
	h.logger.Info(ctx, "new call", zap.String("method", "echo"))

	resp.Message = "echo: " + params.Message

	return nil
}
`, projectName, projectName),
		"app/handler/handler.go": fmt.Sprintf(`package handler

import (
	"context"

	"github.com/khusainnov/%s/app/infra/log"
	"github.com/khusainnov/%s/app/processor"
	reqParams "github.com/khusainnov/%s/specs/server/params/model"
	respParams "github.com/khusainnov/%s/specs/server/response/model"
)

type API interface {
	EchoAPI
}

type EchoAPI interface {
	Echo(ctx context.Context, params *reqParams.EchoParams, resp *respParams.EchoResponse) error
}

// processor interface example
type Processor interface {
	Echo(ctx context.Context) error
}

// Handler example
type Handler struct {
	proc Processor
	logger *log.LoggerCtx
}

func New(logger *log.LoggerCtx, proc *processor.Processor) *Handler {
	return &Handler{
		logger: logger,
		proc: proc,
	}
}
`, projectName, projectName, projectName, projectName),
		"app/infra/log/log_ctx.go": `package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logCtx struct {
	MessageID   string
	HandlerName string
	MsgKey      string
	Message     string
}

type keyType int

const key keyType = iota

func Fields(ctx context.Context, fields []zapcore.Field) []zapcore.Field {
	l, ok := ctx.Value(key).(logCtx)
	if !ok {
		return fields
	}

	if l.MessageID != "" {
		fields = append(fields, zap.String("message_id", l.Message))
	}

	if l.HandlerName != "" {
		fields = append(fields, zap.String("handler_name", l.HandlerName))
	}

	if l.MsgKey != "" {
		fields = append(fields, zap.String("msg_key", l.MsgKey))
	}

	if l.Message != "" {
		fields = append(fields, zap.String("message", l.Message))
	}

	return fields
}

func msg(ctx context.Context, msg string) string {
	if l, ok := ctx.Value(key).(logCtx); ok {
		if l.HandlerName != "" {
			return fmt.Sprintf("%s: %s", l.HandlerName, msg)
		}
	}

	return msg
}

func WithMessageID(ctx context.Context, messageID string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MessageID = messageID

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{MessageID: messageID})
}

func WithHandlerName(ctx context.Context, handlerName string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.HandlerName = handlerName

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{HandlerName: handlerName})
}

func WithMsgKey(ctx context.Context, msgKey string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MsgKey = msgKey

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{MsgKey: msgKey})
}

func WithMessage(ctx context.Context, msg string) context.Context {
	if l, ok := ctx.Value(key).(logCtx); ok {
		l.MsgKey = msg

		return context.WithValue(ctx, key, l)
	}

	return context.WithValue(ctx, key, logCtx{Message: msg})
}

`,
		"app/infra/log/logger.go": `package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Warn(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, fields ...zapcore.Field)
	Fatal(ctx context.Context, msg string, fields ...zapcore.Field)
	Panic(ctx context.Context, msg string, fields ...zapcore.Field)
}

type LoggerCtx struct {
	log *zap.Logger
}

func New(log *zap.Logger) *LoggerCtx {
	return &LoggerCtx{log: log}
}

func (l *LoggerCtx) Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.DebugLevel, msg, fields...)
}

func (l *LoggerCtx) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.InfoLevel, msg, fields...)
}

func (l *LoggerCtx) Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.WarnLevel, msg, fields...)
}

func (l *LoggerCtx) Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.ErrorLevel, msg, fields...)
}

func (l *LoggerCtx) Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.FatalLevel, msg, fields...)
}

func (l *LoggerCtx) Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.logCtx(ctx, zap.PanicLevel, msg, fields...)
}

func (l *LoggerCtx) logCtx(ctx context.Context, lvl zapcore.Level, m string, f ...zapcore.Field) {
	f = Fields(ctx, f)
	l.log.Log(lvl, msg(ctx, m), f...)
}

`,
		"app/infra/system/system.go": fmt.Sprintf(`package system

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	promClient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/khusainnov/%s/app/config"
)

const (
	shutdownTimeout = 5 * time.Second
)

type System struct {
	Registry *promClient.Registry

	meterProvider *sdkmetric.MeterProvider
	prometheus    http.Handler

	Mux *http.ServeMux
	srv *http.Server
}

func New(log *zap.Logger, cfg *config.Config) (*System, error) {
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":53000"
	}

	registry := promClient.NewPedanticRegistry()
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(
			collectors.WithGoCollections(collectors.GoRuntimeMemStatsCollection|collectors.GoRuntimeMetricsCollection),
		),
		collectors.NewBuildInfoCollector(),
	)

	promExporter, err := prometheus.New(
		prometheus.WithRegisterer(registry),
	)
	if err != nil {
		return nil, fmt.Errorf("prometheus: %%w", err)
	}

	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
	)

	mux := http.NewServeMux()

	s := &System{
		Registry: registry,

		meterProvider: metricProvider,
		prometheus:    promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		Mux:           mux,
		srv: &http.Server{
			Handler: mux,
			Addr:    cfg.Server.Addr,
		},
	}

	otel.SetMeterProvider(s.MeterProvider())
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	s.registerRoot()
	s.registerProfiler()
	s.registerPrometheus()

	log.Info("Metrics initialized",
		zap.String("http.addr", s.srv.Addr),
	)

	return s, nil
}

func (s *System) registerProfiler() {
	// Routes for pprof
	s.Mux.HandleFunc("/debug/pprof/", pprof.Index)
	s.Mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.Mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.Mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.Mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Manually add support for paths linked to by index page at /debug/pprof/
	s.Mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	s.Mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	s.Mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	s.Mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func (s *System) registerPrometheus() {
	s.Mux.Handle("/metrics", s.prometheus)
}

func (s *System) MeterProvider() metric.MeterProvider {
	return s.meterProvider
}

func (s *System) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *System) registerRoot() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder
		b.WriteString("Service is up and running.\n\n")
		b.WriteString("\nAvailable debug endpoints:\n")
		for _, s := range []struct {
			Name        string
			Description string
		}{
			{"/metrics", "prometheus metrics"},
			{"/debug/pprof/", "exported pprof"},
		} {
			b.WriteString(fmt.Sprintf("%%-20s - %%s\n", s.Name, s.Description))
		}

		_, _ = fmt.Fprintln(w, b.String())
	})
}

func (s *System) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	wg.Go(func() error {
		// Wait until g ctx canceled, then try to shut down server
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return s.Shutdown(ctx)
	})

	return wg.Wait()
}

`, projectName),
		"app/infra/storage/connect_db.go": fmt.Sprintf(`package infra

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/khusainnov/%s/app/config"
	"gitlab.com/khusainnov/driver/postgres"
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
		return nil, fmt.Errorf("failed to connect to database, %%w", err)
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
	"github.com/khusainnov/%s/app/config"
	"github.com/khusainnov/%s/app/handler"
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
		return fmt.Errorf("failed to register service: %%w", err)
	}

	router := mux.NewRouter()
	router.Handle(HTTP_API_PREFIX, rpcServer)

	s.httpServer = &http.Server{
		Addr:    s.cfg.Addr,
		Handler: router,
	}

	if err := http.ListenAndServe(s.cfg.Addr, router); err != nil {
		return fmt.Errorf("failed to start server: %%w", err)
	}

	return nil
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %%w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown the server: %%w", err)
	}

	return nil
}

`, projectName, projectName),
		"app/processor/processor.go": fmt.Sprintf(`package processor

import (
	"context"

	"github.com/khusainnov/%s/app/infra/log"
)

// Processor example
type Processor struct {
	logCtx *log.LoggerCtx
}

func New(logCtx *log.LoggerCtx) *Processor {
	return &Processor{
		logCtx: logCtx,
	}
}

func (p *Processor) Echo(ctx context.Context) error {
	p.logCtx.Warn(ctx, "implement  me")

	return nil
}

`, projectName),
		"app/service/services.go": fmt.Sprintf(`package service

type Services struct{}

func NewServices() *Services {
	return &Services{}
}
`),
		"app/config/config.go": fmt.Sprintf(`package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	Prefix = "%s"
)

// Config example
type Config struct {
	Server *Server
	DB     *DB
}

type Server struct {
	Addr string `+"`"+`envconfig:"ADDR"`+"`"+`
}

type DB struct {
	Host         string        `+"`"+`envconfig:"HOST"`+"`"+`
	Port         string        `+"`"+`envconfig:"PORT"`+"`"+`
	User         string        `+"`"+`envconfig:"USER"`+"`"+`
	Password     string        `+"`"+`envconfig:"PASSWORD"`+"`"+`
	Name         string        `+"`"+`envconfig:"NAME"`+"`"+`
	SSLMode      string        `+"`"+`envconfig:"SSL_MODE"`+"`"+`
	PingInterval time.Duration `+"`"+`envconfig:"PING_INTERVAL"`+"`"+`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess(Prefix, &c)

	return &c
}
`, strings.ToUpper(projectName)),
		"app/app.go": fmt.Sprintf(`package app

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/khusainnov/%s/app/config"
	"github.com/khusainnov/%s/app/handler"
	"github.com/khusainnov/%s/app/infra/server"
	//"github.com/khusainnov/%s/app/infra/storage" 
	"github.com/khusainnov/%s/app/processor"
	//"github.com/khusainnov/%s/app/repository"
)

// App struct example
type App struct {
	Cfg    *config.Config
	Log    *zap.Logger
	Server *server.Server
	done   chan struct{}
}

func New() *App {
	cfg := config.NewFromEnv()

	app := &App{
		Cfg: cfg,
	}

	log, err := app.CreateLogger()
	if err != nil {
		panic(err)
	}

	app.Log = log

	c := NewContainer(log, cfg)

	app.Log.Info("config init", zap.Any("config", c.Config))

	app.Log.Info("connecting to db")
	/*conn, err := storage.New(app.Log, c.Config.DB)
	if err != nil {
		app.Log.Error("failed to connect to db")
		panic(err)
	}*/

	app.Log.Info("creating new handler")
	handlers := handler.New(
		c.Logger,
		processor.New(c.Logger),
	)

	app.Log.Info("creating new server")
	srv := server.New(c.Config.Server)
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
		"app/container.go": fmt.Sprintf(`package app

import (
	"github.com/khusainnov/%s/app/config"
	"github.com/khusainnov/%s/app/infra/log"
	"github.com/khusainnov/%s/app/service"
	"go.uber.org/zap"
)

type Container struct {
	Config   *config.Config
	Logger   *log.LoggerCtx
	Services *service.Services
	//Prometheus *prometheus.Registry
}

func NewContainer(logger *zap.Logger, cfg *config.Config) *Container {
	return &Container{
		Config:   cfg,
		Logger:   log.New(logger),
		Services: service.NewServices(),
	}
}`, projectName, projectName, projectName),
		"specs/server/params/model/echo.go": `package model

type EchoParams struct {
	Message string ` + "`" + `json:"message"` + "`" + `
}

`,
		"specs/server/response/model/echo.go": `package model

type EchoResponse struct {
	Message string ` + "`" + `json:"message"` + "`" + `
}

`,
		"deployment/local.env": fmt.Sprintf(`# Local environment variables
APP_NAME=MyApp

%s_SERVER_ADDR=:5050
%s_DB_HOST=localhost
%s_DB_PORT=5432
%s_DB_USER=postgres
%s_DB_PASSWORD=postgres
%s_DB_NAME=%sapp
%s_DB_SSL_MODE=disable
%s_DB_PING_INTERVAL=30s
`,
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
			projectName,
			strings.ToUpper(projectName),
			strings.ToUpper(projectName),
		),
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
	"github.com/khusainnov/%s/app"
)

func main() {
	a := app.New()
	a.Run()
}
`, projectName),
		"go.mod": fmt.Sprintf(`module github.com/khusainnov/%s

go 1.24

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/rpc v1.2.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.22.0
	gitlab.com/khusainnov/driver v0.2.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/exporters/prometheus v0.57.0
	go.opentelemetry.io/otel/metric v1.35.0
	go.opentelemetry.io/otel/sdk/metric v1.35.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.14.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

`, projectName),
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

# Installing linter
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.0

# Copy go.mod and go.sum and installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Building binary
RUN go build -o main .

# Stage 2: Run
FROM gcr.io/distroless/base-debian12

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
		"README.md": fmt.Sprintf("# %s\n\nProject description.", projectName),
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
