package main

import (
	"database/sql"
	"flag"
	"jyywong/snippetbox/internal/models"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	db, err := openDB(*dsn)

	app := &application{
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
		snippets: &models.SnippetModel{DB: db},
	}
	if err != nil {
		app.logger.Error(err.Error())
	}
	defer db.Close()

	app.logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
