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


buf.WriteTo(w)
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

### isolating the application routes

The responsibilities of `main()` should be limited to:
- parsing the runtime configuration settings for the application
- establishing the dependencies for the handlers
- running the HTTP server


#### env variables

```go
// returns empty string if no env. variable is provided.
addr := os.Getenv("APP_ADDR")

```
## data-driven responses
database driver: acts as a *middleman* between `MySQL` and our Go application, translating commands between Go and the MySQL database itself.

```sh
# get the latest version with the major release 1
go get github.com/go-sql-driver/mysql@v1

cat go.mod

# is not generated to b
cat go.sum

go mod verify
go mod download
```

upgrading a package
```sh
go get github.com/foo/bar

# vs.
# ⚠️ using `-u` flag increases the risk of breakes when upgrading packages.
go get -u github.com/foo/bar

```

removing unused packages
```sh
go get github.com/foo/bar@none

# automatically removes any unused packages from go.mod and go.sum
go mod tidy

```
📢 Go won't know what "mysql" means until we register `the MySQL driver` with a blank import.
```sh
{"time":"2026-05-25T22:00:11.687718226+02:00","level":"ERROR","msg":"sql: unknown driver \"mysql\" (forgotten import?)"}
exit status 1
```
Don't forget to add this to the imports:
```go
import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

//                       👇
// db, err := sql.Open("mysql", dsn)
```

`sql.Open()` returns a `sql.DB` value, a ppol of many connections. Go manages the connections as needed.

### design the database model
CORE IDEA: 
> create a `standalone models package` (encapsulate the code for working with MySQL in a separate package), so that our database logic is reusable and decoupled from the web application.

Notice that our database logic won't be tied to our handlers. 
A handlre's job is simply validating requests & writing responses.

Go provides 3 different methods for executing database queries:
- `DB.Query()` SELECT
- `DB.QueryRow()` SELECT
- `DB.Exec()`: used for statements not returning rows


## middleware

REMEMBER:
> You can think of a Go web application as a chain of `ServeHTTP()` methods being called one after another.

servemus'x `SErveHTTP()` => handler's `ServeHTTP()`

The pattern for creating our own middleware looks like this:
```go
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// any code here will execute on the way down the chain.

        // stop executing the chain if the user is not authorized.
        if !(isAuthorized(r)) {
            w.WriteHeader(http.StatusForbidden)
            return
        }

		next.ServeHTTP(w, r)
        // any code here will execute on the way back up the chain.
	})
}

```

Where you position the middleware in the chain of handlers will affect the behavior of the application.

In any middleware handler, code which comes before `next.ServeHTTP()` will be executed on the way down the chain, and any code after `next.ServeHTTP()` - or in a deferred function - will be executed on the way back up.

### panic recovery

REMEMBER: Go's HTTP server handles requests concurrently, with each HTTP request handled in its own separate goroutine.

NOTE: Go's HTTP server automatically recovers any panics in the goroutines it created.

> deferred functions in the current goroutine are always called following a panic.

#### panic recovery in background goroutines

> our middleware will only recover panics that happen in the same goroutine that executed the `recoverPanic()` middleware.

### composable middleware chain
install the package
```sh
go get github.com/justinas/alice@v1
```

```go
// converts this:
return mw1(mw2(mw3(myHandler)))

// to:
return alice.New(mw1, mw2, mw3).Then(myHandler)
```

This way, we can assign our middleware chain to variables:
```go

```


## form

```go
r.ParseForm()

r.PostForm.Get("title")
// vs.
r.Form.Get("title")
// vs.
r.FormValue("title")

// ----- query string parameters

r.URL.Query().Get("title")

```

### automatic form parsing

```sh
go get github.com/go-playground/form/v4@v4
```

- initialize a new `*form.Decoder` instance in `main.go` and make it available to the handlers as a dependency
- wrap the application routes with the middleware provided by the `SessionManager.LoadAndServe()` method


## session
install the necessary packages (the session manager)
```sh
go get github.com/alexedwards/scs/v2@v2

go get github.com/alexedwards/scs/mysqlstore@latest
```

`alexewards/scs`:
- server-side only
- automatic loading and saving of session data via middleware
- allows renewal of session IDs
- supports a variety of databases, including MySQL, PostgreSQL and Redis

setting up the session manager, the basics:
- create a session table
- establish a `session manager` in `main.go`, and make it available to the handlers.
- wrap the application routes with the middleware provided by `SessionManager.LoadAndServe()` method
    - do we need the middleware to act on all our application routes? for example on `GET /static/`?
- session manager stores the information temporarily in the **request context**

### behind the scenes of session management
session token: `Cookies` section
The `MySQL BLOB` contains a `gob-encoded` representation of the session data.

The `LoadAndServe()` middleware checks each incoming request for a `session cookie`.
If a session cookie is present, it reads the session token fro the cookie & retrieves the corresponding session data from the database (it checks the expiration data).
Session data is added to the `request context` so it can be used by the http handlers.

## Server & Security Improvements

### http.Server
`http.ListenAndServe()` is very useful in short examples and tutorials.
In real-world applications:
```go
// initialize a new http.Server struct:
srv := &http.Server{
    Addr:     cfg.addr,
    Handler:  app.routes(),
    // force `http.Server` to use our structured logger at Error level for its log messages 
    ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
}

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

### fmt
```go
// s -----
fmt.Sprintf("snippet id %d ...", id)


// f -----
fmt.Fprintf(w, "item id: %d...", id)
// write the snippet data as a plain-text HTTP response body.
fmt.Fprintf(w, "%v", snippet)

// ---
fieldErrors := make(map[string]string)

//...
if len(fieldErrors) > 0 {
    fmt.Fprint(w, fieldErrors)
    return
}
// ---



```

### errors

```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return Snippet{}, ErrNoRecord
    } else {
        return Snippet{}, err
    }
}

var ErrNoRecord = errors.New("models: no matching record found")

# -----

if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
}

# -----

```


### Status Code
```yml
# 2xx
200     OK
201     created

# 3xx
301     permanent redirect
302     temporarily redirect
303     see other

# 4xx
400     bad request
404     not found
405     method not allowed
422     unprocessable content/entity

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

http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

```

### curl
```sh
curl -i localhost:8080/

curl --head localhost:8080/
curl -I localhost:8080/

curl -i -d "" localhost:8080/
curl -i -X POST localhost:8080/

# -L: automatically follow redirect
curl -iL -d "" localhost:8080/snippet/create
```

### str

```go
utf8.RuneCountInString(title)

strings.TrimSpace(title) == ""



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
## RND

### TCP

A TCP connection is uniquely identified by:
> (source IP, source port, destination IP, destination port)


## advanced

### map vs. struct
Unlike `struct` fields, `map` key names don't have to be capitalized in order to access them from a template.

### val vs. ref
```go

type templateData struct {
	Snippet models.Snippet
}

```
What this means? vs. `Snippet *models.Snippet`
- Snippet is always present
- You get a copy
- No nil possible

### defer
This means:


### map

```go
// initialize
fieldErrors := make(map[string]string)


// vs.


// -----

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) AddFieldError(key, msg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}


// -----

```

```go
defer func() {
    ...
}()
```
> Defer the execution of this anonymous function call.

Best equivalent in python is
```py
def do_work():
    try:
        print("doing work")
        raise ValueError("boom")
    finally:
        print("cleanup runs no matter what")
```
or
```py
with open("data.txt") as f:
    data = f.read()

# f = open("data.txt")
# try:
#     data = f.read()
# finally:
#     f.close()
```


### generic function
Generic functions work with values of different types.
```go
func PermittedValue[T comparable](val T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, val)
}


// -----

func count[T comparable](v T, s []T) int {
    count := 0
    for _, itm := range s {
        if v == itm {
            count++
        }
    }
    return count
}

```
## WTF

2.10
3.3     closures for dependency injection
5.6     custom template functions
6.2     CSP (content-security policy) section
6.4     panic recovery
