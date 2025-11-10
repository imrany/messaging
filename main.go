package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	v1 "github.com/imrany/whats-email/internal/v1"
	customMiddleware "github.com/imrany/whats-email/middleware"
	"github.com/imrany/whats-email/pkg/whatsapp"

	_ "modernc.org/sqlite"
)

func createServer() *http.Server {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// when the request has timed out and further processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(customMiddleware.CorsMiddleware)

	// Public routes
	r.Get("/health", v1.HealthHandler)

	// Protected routes
	r.Route("/api/v1", func(r chi.Router) {
		// r.Use(middleware.AuthMiddleware) // Add your authentication middleware here

		r.Post("/mailer/send", v1.SendMail)
		r.Post("/whatsapp/send", v1.SendWhatsAppMessage)
	})

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", viper.GetString("HOST"), viper.GetInt("PORT")),
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	return srv
}

func runServer() {
	var err error
	server := createServer()
	port := viper.GetInt("PORT")
	host := viper.GetString("HOST")

	// Initialize WhatsApp client
	slog.Info("Initializing WhatsApp client...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := whatsapp.Init(ctx, nil); err != nil {
		slog.Error("Error initializing WhatsApp client", "error", err.Error())
		slog.Warn("Server will start without WhatsApp integration")
		// Continue running server even if WhatsApp fails to initialize
	} else {
		slog.Info("WhatsApp client initialized successfully")
	}

	// Start server in goroutine
	go func() {
		slog.Info("Server started", "host", host, "port", port)
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "error", err.Error())
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	slog.Info("Shutdown signal received", "signal", sig)

	// Shutdown WhatsApp client
	slog.Info("Disconnecting WhatsApp client...")
	whatsapp.Disconnect()

	// Shutdown HTTP server
	slog.Info("Shutting down HTTP server...")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	} else {
		slog.Info("Server exited cleanly")
	}
}

func main() {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using defaults")
	} else {
		slog.Info(".env file loaded successfully")
	}

	// root command
	var rootCmd = &cobra.Command{
		Use:   "smart-spore-hub",
		Short: "Smart Spore Hub",
		Long:  "Smart Spore Hub is a web application for managing spore data.",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}

	// flags
	rootCmd.PersistentFlags().Int("port", 8080, "Port to listen on")
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "Host to listen on")
	rootCmd.PersistentFlags().String("SMTP_HOST", "smtp.gmail.com", "SMTP HOST (env: SMTP_HOST)")
	rootCmd.PersistentFlags().Int("SMTP_PORT", 587, "SMTP PORT (env: SMTP_PORT)")
	rootCmd.PersistentFlags().String("SMTP_USERNAME", "", "SMTP Username (env: SMTP_USERNAME)")
	rootCmd.PersistentFlags().String("SMTP_PASSWORD", "", "SMTP Password (env: SMTP_PASSWORD)")
	rootCmd.PersistentFlags().String("SMTP_EMAIL", "", "SMTP Email (env: SMTP_EMAIL)")

	// Bind flags to viper
	viper.BindPFlag("PORT", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("HOST", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("SMTP_HOST", rootCmd.PersistentFlags().Lookup("SMTP_HOST"))
	viper.BindPFlag("SMTP_PORT", rootCmd.PersistentFlags().Lookup("SMTP_PORT"))
	viper.BindPFlag("SMTP_USERNAME", rootCmd.PersistentFlags().Lookup("SMTP_USERNAME"))
	viper.BindPFlag("SMTP_PASSWORD", rootCmd.PersistentFlags().Lookup("SMTP_PASSWORD"))
	viper.BindPFlag("SMTP_EMAIL", rootCmd.PersistentFlags().Lookup("SMTP_EMAIL"))

	// Bind env variables
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
