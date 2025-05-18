package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
)

type config struct {
	port int
	env  string
	mail struct {
		host       string
		port       string
		user       string
		pwd        string
		recipients string
		allowed_ip string
	}
}

type envelope map[string]interface{}

type application struct {
	config    config
	wg        sync.WaitGroup
	validator *validator.Validate
}

func init() {

	var programLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
}

func main() {

	var cfg config

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		slog.Error("failed to parse port", "error", err)
		os.Exit(1)
	}

	flag.IntVar(&cfg.port, "port", port, "Port to listen on")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Environment (development|production)")
	flag.StringVar(&cfg.mail.host, "MAIL HOST", os.Getenv("EMAIL_HOST"), "MAIL HOST")
	flag.StringVar(&cfg.mail.port, "MAIL PORT", os.Getenv("EMAIL_PORT"), "MAIL PORT")
	flag.StringVar(&cfg.mail.user, "MAIL USER ", os.Getenv("EMAIL_USER"), "MAIL USER")
	flag.StringVar(&cfg.mail.pwd, "MAIL PASSWORD", os.Getenv("EMAIL_PASS"), "MAIL PWD")
	flag.StringVar(&cfg.mail.recipients, "RECEPIENTS", os.Getenv("RECEPIENTS"), "RECEPIENTS")
	flag.StringVar(&cfg.mail.allowed_ip, "ALLOWED_IP", os.Getenv("ALLOWED_IP"), "ALLOWED_IP")

	flag.Parse()

	app := &application{
		config:    cfg,
		validator: validator.New(),
	}

	err = app.serve()
	if err != nil {
		slog.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
