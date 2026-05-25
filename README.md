# SBox

```go
println("Hello, World!")
```

## Foundation

```sh
go version

```

Think of `module path` as the identifier for your project.

```sh
cd sbox

# turn your project directory into a module.
go mod init github.com/gesatessa/sbox

```
`go.mod` will have your `go version` for the project.

### web app basics

3 essentials
- `handler`: like controllers in MVC pattern. responsible for execuring application logic & writing HTTP response headers and bodies. 

- `servemux`: the router.

- `web server`: In Go, you don't need an external 3rd party server like Nginx or Caddy.


#### w & r
`http.ResponseWriter`: provides methods for assembling an HTTP response & sending it to the user.

`*http.Request`: is a pointer to a struct, holding information about the current request.
- r.Method
- r.URL.Path

```sh
# /tmp/ & then run
go run .
```

### routing requests

`subtree path pattern`: "/" or "/static/"

Additional servemux features:
- Request URL paths are automatically sanitized.
- if `/foo/` is registered, requests `/foo` will be redirected to `/foo`, with `301 Permanent Redirect`


Precedence & conflicts:
The most specific route pattern wins.

### customizing responses

`w.Write()` sends a 200 status code by default. => `w.WriteHeader()`

#### customizing header

```go
w.Header().Add("Server", "Go")

```

#### writing response bodies

...because the `http.ResponseWriter` value in our handlers has a `Write()` method, it satisfies the `io.Writer` interface.

At a practical level, this means that any functions where you see an `io.Writer` parameter, you can pass in the `http.ResponseWriter` value and whatever is being written will subsequently be sent as the body of the HTTP response.
```go
w.Write([]byte("Hellooo))

// 👇
io.WriteString(w, "Hellooo")
fmt.Fprint(w, "Helloooo")

```

content sniffing:
```go
w.Header().Set("Content-Type", "application/json")

```

### project structure
Do NOT over-complicate things.
Try hard to only add structure and complexity when it's demonstrably needed.

- cmd (application-specific code)
cmd/web
cmd/cli (to automate some administrative tasks)
- internal (ancillary non-application-specific code)
- ui

### html templating & inheritance

```go
ts, err := template.ParseFiles(files...)

// base: template name (invokes other templates)
err = ts.ExecuteTemplate(w, "base", nil)

```

{{define "base"}}...{{end}}
{{template "title" .}}

### serving static files

```sh
mkdir -p ./ui/static
curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/
```

```yml
ui/static
├── css
│   └── main.css
├── img
│   ├── favicon.ico
│   └── logo.png
└── js
    └── main.js
```

`http.FileServer` - from `net/http/` package - serves files over HTTP from a specific directory.

```go
fileServer := http.FileServer(http.Dir("./ui/static/"))
http.Handle("/static/", http.StripPrefix("/static", fileServer))

```

Don't forget to update the `base` template.

`http.FileServer` supports "range requests", enabling resumable downloads.
```sh
curl -i -H "Range: bytes=100-199" --output - http://localhost:8080/static/img/logo.png
# HTTP/1.1 206 Partial Content
# Accept-Ranges: bytes
# Content-Length: 100
# Content-Range: bytes 100-199/1075
# Content-Type: image/png
# Last-Modified: Thu, 04 May 2017 13:07:52 GMT
# Date: Mon, 25 May 2026 09:20:53 GMT

# h�j��ZbK�&�"b��dS�"V��M�PQ�S�T��x�PMC1���&�.(غ� ����&�"^"� ZI% 
```

for frequently-served files, at least, it's highly likely that `http.FileServer` will be serving static files in `./ui/static` from RAM rather than making the relatively slow round-trip to hard disk.

## configuration & error handling

### managing configuration settings
- separation between configuration settings & code
- managing configuration settings at runtime (based on the environment: dev, testing, prod)

```sh
go run ./cmd/web -help
# Usage of /home/skye/.cache/go-build/8e/8e95...dea-d/web:
#   -addr string
#         HTTP network address (default ":8080")

go run ./cmd/web -addr=":4000
```

### structured logging
the `log/slog` package from the standard library.

NOTE: all structured loggers have a structured logging handler associated with them (not to be confused with an HTTP handler).

```go
addr := flag.String("addr", ":8080", "HTTP network address")

// nil: no customization
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.Info("starting server", "addr", *addr)
// {"time":"2026-05-25T13:21:04.231","level":"INFO","msg":"starting server","addr":":8080"}
```

// Debug, Info, Warn, Error

In `staging` or `production` environments, we can redirect the log stream for archival.
IMPORTANT: the final destination of the logs can be managed by the execution environment independently of the application.
```go
// redirect the standard out stream to an on-disk file
go run ./cmd/web >> /tmp/web.log

```

NOTE: custom loggers created by `slog.New()` are concurrency-safe.
We can share a single logger and use it across multiple goroutines and in the HTTP handlers without needing to worry about race conditions.
That said, ...

### dependency injection

Most web applications will have multiple dependencies that their handlers need to access:
- db connection pool,
- centralized error handlers
> Q: How can we make any dependency available to our handlers?

=> **inject dependencies into the handlers**
This - in comparison to having a gloval variable - has the benefit to:
- make the code more explicit
- less error-prone
- easier to unit test

> put the dependencies into a custom `application` struct, and define the handlers as method against it.

#### env variables

```go
// returns empty string if no env. variable is provided.
addr := os.Getenv("APP_ADDR")

```

## MiSK

### r

```go
r.PathValue("id") // always a string

```

### log
```go
log.Print("starting server on :8080")

log.Print("starting server on ", *addr)
log.Printf("starting server on %s", *addr)

log.Fatal(err)
```

# fmt ----
```go
// s
fmt.Sprintf("snippet id %d ...", id)

// f
fmt.Fprintf(w, "item id: %d...", id)
```

### Status Code
```yml
# 2xx
200     OK
201     created

# 3xx
301     permanent redirect
302     temporarily redirect

# 4xx
400     bad request
404     not found
405     method not allowed

# 5xx
500     internal server error
502

503
```


### http
```go
http.NotFound(w, r)

if err != nil {
    log.Print(err.Error())
    // log.Printf("error parsing template: %v", err)

    http.Error(w, "internal server error", 5xx)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
}

```

### curl
```sh
curl -i localhost:8080/

curl --head localhost:8080/
curl -I localhost:8080/

curl -i -d "" localhost:8080/
curl -i -X POST localhost:8080/
```

## project structure

```yml
# cmd
./cmd/web 
- main.go   - handlers.go


# internal



# ui
./ui/html

./ui/static
```

## WTF

2.10
3.3     closures for dependency injection