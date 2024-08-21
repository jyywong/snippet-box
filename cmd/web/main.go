package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	app.logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)
	app.logger.Error(err.Error())
	os.Exit(1)
}
