cdncache
========

[![GoDoc](https://godoc.org/github.com/go-on/cdncache?status.png)](http://godoc.org/github.com/go-on/cdncache)

a cache for CDN files

Example
--------


```go
package main

import (
    "fmt"
    "gopkg.in/go-on/cdncache.v1"
    "net/http"
)

// mounts the cached files at /cdn-cache/
var cdn1 = cdncache.CDN("/cdn-cache/")

// does no caching (empty mountpoint)
var cdn2 = cdncache.CDN("")

func serve(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,
        `<html>
        <body>
        <a href="%s">local (cached)</a><br>
        <a href="%s">remote (CDN)</a><br>
        </body>
        </html>`,
        cdn1("//code.jquery.com/jquery-1.11.0.min.js"),
        cdn2("//code.jquery.com/jquery-1.11.0.min.js"),
    )
}

func main() {
    http.Handle("/", http.HandlerFunc(serve))
    http.ListenAndServe(":8383", nil)
}
```
