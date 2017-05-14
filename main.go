package main

import (
	"log"
	"net/http"
)

const (
	VERSION = "0.2.0"
	AUTHOR  = "交大吴彦祖"
)

func main() {
	http.HandleFunc("/user/", userController)
	http.HandleFunc("/items/", itemController)
	http.HandleFunc("/hello", helloController)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
