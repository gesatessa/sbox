package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// hold the application configuration settings.
type config struct {
	addr   string
	static string
}

// hold the application-wide dependencies.
type application struct {
	logger *slog.Logger
}

func main() {
	var cfg config

	// flag.String returns a pointer to a string variable that stores the value of the command-line flag.
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.static, "static", "./ui/static/", "path to static files")

	// if any errors occur during flag parsing, (e.g., flag value cannot be converted to the expected type),
	// the program will print an error message and exit with a non-zero status code.
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	// initialize a new instance of the application struct, containing the dependencies.
	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.static))
	// register the file server as the handler for all URL paths starting with "/static/".
	// The http.StripPrefix function is used to remove the "/static" prefix from the URL path before the file server handles the request.
	// For example, if the URL path is "/static/css/style.css",
	// the file server will look for the file "./ui/static/css/style.css" on the filesystem.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register the home function as the handler for the root path.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)

	logger.Info("starting server", "addr", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	// any error returned by ListenAndServe is ALWAYS non-nil. Log the error and exit the program.
	// logger.Error("failed to start server", "error", err)
	logger.Error(err.Error())
	os.Exit(1)

}
