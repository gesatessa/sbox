package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// register the file server as the handler for all URL paths starting with "/static/".
	// The http.StripPrefix function is used to remove the "/static" prefix from the URL path before the file server handles the request.
	// For example, if the URL path is "/static/css/style.css",
	// the file server will look for the file "./ui/static/css/style.css" on the filesystem.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register the home function as the handler for the root path.
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)

	log.Print("starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	// any error returned by ListenAndServe is ALWAYS non-nil. Log the error and exit the program.
	log.Fatal(err)

}
