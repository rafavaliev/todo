package main

import (
	"context"
	"net/http"
	"time"
	http2 "todo/handler/http"
	internalDB "todo/internal/db"
	internalLog "todo/internal/log"
	"todo/internal/server"
	"todo/search"
	"todo/task"
	"todo/user"
)

func main() {

	logger := internalLog.New()

	logger.Info("Starting the app")

	db, err := internalDB.New()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, user.UserContextKey, user.User{ID: 1})

	// Migrate the schema
	_ = db.AutoMigrate(&search.SQLUserIndex{}, &task.Task{}, &user.User{})

	searchRepo := search.NewSQLRepository(db)
	taskRepo := task.NewSQLRepository(db)
	userRepo := user.NewSQLRepository(db)

	searchService := search.NewService(searchRepo)
	taskService := task.NewService(taskRepo, searchService)

	srv := server.New(http2.NewHandler(logger, taskService, searchService, userRepo))
	logger.With("addr", srv.Addr).Info("Starting the server")

	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		<-ctx.Done()

		logger.Info("Stopping the server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.With("error", err).Fatal("Could not gracefully shutdown the server.")
		}
		logger.Info("Server stopped.")
		close(done)
	}(done)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("Server stopped")
	}
	<-done

}
