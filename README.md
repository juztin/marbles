Marbles
=======

----

Marbles is a simple collection of tools for web tasks.

 * [Listeners](#Listeners)
 * [Encoders](#Encoders)
 * [Routes](#Routes)

----

## Listeners ##

**HTTP listener**

```
func main() {
    l, err := listeners.NewHTTP("", "9000")
    http.Serve(l, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.responseWriter) {
    fmt.Println("hola")
}
```

**HTTPS listener**

```
func main() {
    l, err := listeners.NewTLS("", "9000", "./server.crt", "./server.key")
    http.Serve(l, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.responseWriter) {
    fmt.Println("hola")
}
```

**Unix socket**

```
func main() {
    l, err := listeners.NewSOCK("/var/run/dirty.sock", os.ModePerm)
    http.Serve(l, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.responseWriter) {
    fmt.Println("hola")
}
```

----

## Encoders ##

**jsonxml** - is currently the only encoder available.

```
func handler(w http.ResponseWriter, r *http.Request) {
    b := []string{"one", "two"}
    jsonxml.Write(w, r, b)
    // Or to also set the StatusCode
    // jsonxml.Write(w, r, b, 201)
}
```

The encoder looks at the `Accept` header to help determine the encoder to use.  

`application/json` is the default if there is no `Accept` header supplied.  

If the content type is unset, or `application/javascript` and there is a `callback` query-string parameter. The resulting JSON data will be wrapped in padding.

## Routes ##

Just some simple routing helpers.  
*The main difference is that these break the net/http `Handler` interface. The `http.ResponseWriter` and `http.Request` get wrapped within a `Context` struct.*  

*(I will be deprecating these in the future. Either getting rid of them completely, and just using something like [Gorilla][1], or updating them to follow the `Handler` interface)*

[1]: http://www.gorillatoolkit.org/