package main

import (
    "net/http"
    "log"
)

func main() {
    http.Handle("/", http.FileServer(http.Dir("./static")))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
