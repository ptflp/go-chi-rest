package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ptflp/go-chi-rest/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ptflp/go-chi-rest/database"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	servePort := os.Getenv("SERVE_PORT")

	db, err := database.Connect(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		fmt.Println("database:", err)
		os.Exit(-1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	postHandlers := handlers.NewPostHandler(db)
	r.Route("/", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandlers.Fetch)
			r.Get("/{id:[0-9]+}", postHandlers.GetByID)
			r.Post("/", postHandlers.Create)
			r.Put("/{id:[0-9]+}", postHandlers.Update)
			r.Delete("/{id:[0-9]+}", postHandlers.Delete)
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", servePort),
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen and serve error", err)
		}
	}()

	fmt.Println("Server started at port:", servePort)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("shutdown error", err)
	}
}
