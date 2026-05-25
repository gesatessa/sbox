package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// hold the application configuration settings.
type config struct {
	addr   string
	static string
	dsn    string
}

// hold the application-wide dependencies.
type application struct {
	logger *slog.Logger
	cfg    config
}

func main() {
	var cfg config

	// flag.String returns a pointer to a string variable that stores the value of the command-line flag.
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.static, "static", "./ui/static/", "path to static files")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DSN"), "MySQL DSN (Data Source Name)")

	// if any errors occur during flag parsing, (e.g., flag value cannot be converted to the expected type),
	// the program will print an error message and exit with a non-zero status code.
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// initialize a new instance of the application struct, containing the dependencies.
	app := &application{
		logger: logger,
		cfg:    cfg,
	}

	logger.Info("starting server", "addr", cfg.addr)
	err = http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB() is a helper function which returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
