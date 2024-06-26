package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/citadel-corp/cats-social/internal/cat"
	catmatch "github.com/citadel-corp/cats-social/internal/cat_match"
	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/citadel-corp/cats-social/internal/common/middleware"
	"github.com/citadel-corp/cats-social/internal/user"
	"github.com/gorilla/mux"
	"github.com/lmittmann/tint"
)

func main() {
	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339,
	})
	slog.SetDefault(slog.New(slogHandler))

	// Connect to database
	// env := os.Getenv("ENV")
	// sslMode := "disable"
	// if env == "production" {
	// 	sslMode = "verify-full sslrootcert=ap-southeast-1-bundle.pem"
	// }
	// connStr := "postgres://[user]:[password]@[neon_hostname]/[dbname]?sslmode=require"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARAMS"))
	// dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), sslMode)
	db, err := db.Connect(connStr)
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize cat domain
	catRepository := cat.NewRepository(db)
	catService := cat.NewService(catRepository)
	catHandler := cat.NewHandler(catService)

	// initialize cat match domain
	catMatchRepository := catmatch.NewRepository(db)
	catMatchService := catmatch.NewService(catMatchRepository, catRepository)
	catMatchHandler := catmatch.NewHandler(catMatchService)

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecoverer)
	v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		io.WriteString(w, "Service ready")
	})

	// user routes
	ur := v1.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/register", userHandler.CreateUser).Methods(http.MethodPost)
	ur.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)

	// cat match routes
	cmr := v1.PathPrefix("/cat/match").Subrouter()
	cmr.HandleFunc("", middleware.Authorized(catMatchHandler.GetCatMatchList)).Methods(http.MethodGet)
	cmr.HandleFunc("", middleware.Authorized(catMatchHandler.Create)).Methods(http.MethodPost)
	cmr.HandleFunc("/approve", middleware.Authorized(catMatchHandler.Approve)).Methods(http.MethodPost)
	cmr.HandleFunc("/reject", middleware.Authorized(catMatchHandler.Reject)).Methods(http.MethodPost)
	cmr.HandleFunc("/{id}", middleware.Authorized(catMatchHandler.Delete)).Methods(http.MethodDelete)

	// cat management routes
	cr := v1.PathPrefix("/cat").Subrouter()
	cr.HandleFunc("", middleware.Authorized(catHandler.GetCatList)).Methods(http.MethodGet)
	cr.HandleFunc("", middleware.Authorized(catHandler.CreateCat)).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", middleware.Authorized(catHandler.UpdateCat)).Methods(http.MethodPut)
	cr.HandleFunc("/{id}", middleware.Authorized(catHandler.DeleteCat)).Methods(http.MethodDelete)

	httpServer := &http.Server{
		Addr:     ":8080",
		Handler:  r,
		ErrorLog: slog.NewLogLogger(slogHandler, slog.LevelError),
	}

	go func() {
		slog.Info(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("HTTP server error: %v", err))
		}
		slog.Info("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	slog.Info(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error: %v", err)
	}
	slog.Info("Shutdown complete.")
}
