package main

import (
    "net/http"
    "io"
)

func helloController(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "I love U!\n")
}
