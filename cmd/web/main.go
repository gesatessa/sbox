package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	// we need the driver's init() function to be called to register the MySQL driver with the database/sql package.
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/form/v4"

	"github.com/gesatessa/sbox/internal/models"
)

// hold the application configuration settings.
type config struct {
	addr   string
	static string
	dsn    string
}

// hold the application-wide dependencies.
type application struct {
	logger         *slog.Logger
	cfg            config
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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
	// make sure the database connection pool is closed before the main() function exits.
	// This will help to prevent resource leaks
	// and ensure that all database connections are properly released when the application shuts down.
	defer db.Close()

	// initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// initialize a decoder instance to be added to the application dependencies.
	formDecoder := form.NewDecoder()

	// initialize & configure a new session manager:
	sessionManager := scs.New() // returns a pointer to the SessionManager struct.
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = time.Hour * 1

	// initialize a new instance of the application struct, containing the dependencies.
	app := &application{
		logger:         logger,
		cfg:            cfg,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

	// Ping() is used to verify that the database connection is alive and working properly.
	// It sends a simple query to the database and waits for a response.
	// If the connection is successful, it returns nil. If there is an error (e.g., network issue, authentication failure, etc.),
	// it returns an error value describing the problem.
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
