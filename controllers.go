package main

import (
    "net/http"
    "io"
    "strings"
)

func helloController(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "Hello world!\n")
}

func userController(w http.ResponseWriter, req *http.Request) {
    path := req.URL.EscapedPath()
    method := strings.Split(path, "/")[2]

    switch method {
        case "login" :
            loginUserController(w, req)
        case "logout" :
            logoutUserController(w, req)
        case "info" :
            infoUserController(w, req)
        case "register" :
            registerUserController(w, req)
        default:
            NotFound(w, req)
    }
}

func itemController(w http.ResponseWriter, req *http.Request) {
    path := req.URL.EscapedPath()
    method := strings.Split(path, "/")[2]

    switch method {
        case "add" :
            addItemController(w http.ResponseWriter, req *http.Request)
        case "list" :
            listItemController(w http.ResponseWriter, req *http.Request)
        case "share" :
            shareItemController(w http.ResponseWriter, req *http.Request)
        default:
            NotFound(w, req)
    }
}

func registerUserController(w http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        username := req.FormValue("username")
        passwd := req.FormValue("password")
        email := req.FormValue("email")
        register(username, passwd, email)
    }
}
