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



## MiSK

### r

```go
r.PathValue("id") // always a string

```

### log
```go
log.Print("starting server on :8080")

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
