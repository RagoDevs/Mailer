package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		slog.Info("shutting down gracefully, press Ctrl+C again to force")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.LogAttrs(context.Background(),
			slog.LevelInfo,
			"completing background tasks",
			slog.String("addr", srv.Addr),
		)

		app.wg.Wait()
		shutdownError <- nil

	}()

	slog.LogAttrs(context.Background(),
		slog.LevelInfo,
		"Starting server ",
		slog.String("addr", srv.Addr),
		slog.String("env", app.config.env))

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.LogAttrs(context.Background(),
		slog.LevelInfo,
		"stopped server",
		slog.String("addr", srv.Addr),
	)

	return nil

}
