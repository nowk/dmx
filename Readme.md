# dmx

[![Build Status](https://travis-ci.org/nowk/dmx.svg?branch=master)](https://travis-ci.org/nowk/dmx)
[![GoDoc](https://godoc.org/github.com/nowk/dmx?status.svg)](http://godoc.org/github.com/nowk/dmx)

A simple pattern match mux. *A speed experiment.*


## Install

    go get gopkg.in/nowk/dmx.v2


## Example

    package main

    import "net/http"
    import "gopkg.in/nowk/dmx.v2"

    var getPostHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")
    })

    func main() {
      mux := dmx.New()
      mux.Get("/posts/:id", getPostHandler)

      h := mux.Handler(dmx.NotFound(mux))

      err := http.ListenAndServe(":3000", h)
      if err != nil {
        log.Fatalf("fatal: listen: %s", err)
      }
    }

##### Basic Method Shortcuts

GET

    mux.Get(string, http.Handler)
    
POST
    
    mux.Post(string, http.Handler)
    
PUT

    mux.Put(string, http.Handler)
    
DELETE

    mux.Del(string, http.Handler)

---

Use the `Add` method for binding multiple methods to a path and handler.

    mux.Add("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // ...
    }), "PUT", "PATCH")

---

`.:format` at the end of your path pattern will parse the extension and provide it as a parameter named `:_format`.

    mux.Get("/posts/:id.:format", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      v := req.URL.Query()
      id := v.Get(":id")
      format := v.Get(":_format")

      // ...
    }))

Using `.:format` will match paths with or without the extension.

---

##### Handling Not Found

You define how you want non-matches to be handled. 

You can use the built in `dmx.NotFound`. This will return either `404` or `405` (with an `Allow` header).

    h := mux.Handler(dmx.NotFound(mux))


Or you you can pass it off to another handler. Ex: a file server

    stc := http.Dir("./public")
    pub := http.FileServer(stc)
    ...
    h := mux.Handler(pub)


## License

MIT

