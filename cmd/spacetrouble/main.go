package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonzhukov/spacetrouble/internal/handler"

	"github.com/apex/log"

	"github.com/antonzhukov/spacetrouble/internal"
	"github.com/antonzhukov/spacetrouble/internal/launch"

	"github.com/antonzhukov/spacetrouble/internal/booking"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/lib/pq"
)

func main() {
	cfg := initConfig()
	logger := initLogger()
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database,
	)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Fatal("sql.Open failed", zap.Error(err), zap.String("conn", connString))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("db.Ping failed", zap.Error(err), zap.String("conn", connString))
	}
	store := booking.NewPostgres(db)
	launches := launch.NewSpaceX(
		&http.Client{Timeout: 10 * time.Second},
		"https://api.spacexdata.com/v4/launches",
	)
	bookingCenter := internal.NewBookingCenter(logger, launches, store)
	bookingHandler := handler.NewBooking(logger, bookingCenter)

	r := mux.NewRouter()
	r.HandleFunc("/bookings", bookingHandler.AddBooking).Methods("POST")
	r.HandleFunc("/bookings", bookingHandler.GetBookings).Methods("GET")

	srv := http.Server{
		Addr:         ":8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			logger.Fatal("srv.ListenAndServe failed", zap.Error(err))
		}
	}()
	logger.Info("Server started", zap.Int("port", 8000))

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, os.Interrupt)
	sig := <-c

	// Shutdown http server
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Fatal("HTTP server shutting down")
	}

	logger.Info("application shutting down...", zap.String("signal", sig.String()))
	os.Exit(1)
}

type Config struct {
	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
}

func initConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Postgres.Host, "pg.host", "", "Postgres host")
	flag.IntVar(&cfg.Postgres.Port, "pg.port", 5432, "Postgres port")
	flag.StringVar(&cfg.Postgres.User, "pg.user", "", "Postgres username")
	flag.StringVar(&cfg.Postgres.Password, "pg.password", "", "Postgres password")
	flag.StringVar(&cfg.Postgres.Database, "pg.db", "", "Postgres database")
	flag.Parse()
	return cfg
}

func initLogger() *zap.Logger {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encCfg
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("createLogger failed: %v", err)
		os.Exit(1)
	}
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	return logger
}
