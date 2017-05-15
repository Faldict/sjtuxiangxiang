package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func helloController(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello world!\n")
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "No find!\n")
}

func userController(w http.ResponseWriter, req *http.Request) {
	path := req.URL.EscapedPath()
	method := strings.Split(path, "/")[2]

	switch method {
	case "login":
		loginUserController(w, req)
	case "logout":
		logoutUserController(w, req)
	case "info":
		infoUserController(w, req)
	case "register":
		registerUserController(w, req)
	default:
		NotFound(w, req)
	}
}

func itemController(w http.ResponseWriter, req *http.Request) {
	path := req.URL.EscapedPath()
	method := strings.Split(path, "/")[2]

	switch method {
	case "add":
		addItemController(w, req)
	case "list":
		listItemController(w, req)
	case "share":
		shareItemController(w, req)
	default:
		NotFound(w, req)
	}
}

func registerUserController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		username := req.FormValue("username")
		passwd := req.FormValue("password")
		email := req.FormValue("email")
		rst := register(username, passwd, email)
		w.Write(rst)
	}
}

func loginUserController(w http.ResponseWriter, req *http.Request) {
	var rst []byte
	if req.Method == "POST" {
		username := req.FormValue("username")
		passwd := req.FormValue("password")
		//rst := login(w, username, passwd) // 不清楚w和req能不能这样传递
		db, err := sql.Open("mysql", "user:password@/dbname")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer db.Close()

		var PSD string
		err = db.QueryRow("SELECT user_PSD FROM users WHERE user_ID=?", username).Scan(&PSD)
		if err != nil {
			log.Fatal(err.Error())
			rst = []byte("300002")
			goto Here
		}

		if passwd == PSD {
			cookie := http.Cookie{
				Name:   "uid",
				Value:  username,
				Path:   "/",
				MaxAge: 86400,
			}
			http.SetCookie(w, &cookie)
			rst = []byte("200000")
			goto Here
		} else {
			rst = []byte("300001")
			goto Here
		}
	Here:
		w.Write(rst)
	}
}

func logoutUserController(w http.ResponseWriter, req *http.Request) {
	//rst := logout(w, *req) 感觉有问题不同这个了
	/*
		uid, err := req.Cookie("uid")
		if err != nil {
			log.Fatal(err.Error())
			return []byte("Cookie read error")
		}
	*/
	var uid string
	for _, cookie := range req.Cookies() {
		uid = cookie.Name
		break
	}

	cookie := http.Cookie{
		Name:   "uid",
		Value:  uid,
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	rst := []byte("Logout successfully")
	w.Write(rst)
}

func infoUserController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		uid := req.FormValue("uid")
		info := userinfo(uid)
		w.Write(info)
	}
}

func addItemController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		uid := req.FormValue("uid")
		obj_name := req.FormValue("obj_name")
		obj_price := req.FormValue("obj_price")
		obj_info := req.FormValue("obj_info")
		use_time := req.FormValue("use_time")
		rst := addItem(obj_name, uid, obj_price, obj_info, use_time)
		w.Write(rst)
	}
}

func listItemController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		rst := listItem()
		w.Write(rst)
	}
}

func shareItemController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		uid := req.FormValue("uid")
		obj_name := req.FormValue("obj_name")
		obj_price := req.FormValue("obj_price")
		obj_info := req.FormValue("obj_info")
		use_time := req.FormValue("use_time")
		rst := shareItem(obj_name, uid, obj_price, obj_info, use_time)
		w.Write(rst)
	}
}
