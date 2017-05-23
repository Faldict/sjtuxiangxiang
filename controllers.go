package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
	io.WriteString(w, "404 Not Found")
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

func messageController(w http.ResponseWriter, req *http.Request) {
	path := req.URL.EscapedPath()
	method := strings.Split(path, "/")[2]

	switch method {
	case "send":
		sendMessageController(w, req)
	case "receive":
		receiveMessageController(w, req)
	case "listMessage":
		listMessageController(w, req)
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
	case "listItem":
		listItemController(w, req)
	case "shareRequest":
		shareRequestController(w, req)
	case "shareResponse":
		shareResponseController(w, req)
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
	var PSD string
	if req.Method == "POST" {
		username := req.FormValue("username")
		passwd := req.FormValue("password")
		db, err := sql.Open("mysql", "user:password@/dbname")
		if err != nil {
			log.Fatal(err.Error())
			rst = []byte("300001")
			goto Here
		}
		defer db.Close()

		err = db.QueryRow("SELECT user_PSD FROM users WHERE user_ID=?", username).Scan(&PSD)
		if err != nil {
			log.Fatal(err.Error())
			rst = []byte("300004") //300002SELECT错误
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
			rst = []byte("200000") //200000登录成功
			goto Here
		} else {
			rst = []byte("200001") //300001密码错误
			goto Here
		}
	Here:
		w.Write(rst)
	}
}

func logoutUserController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value
	cookie := http.Cookie{
		Name:   "uid",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	rst := []byte(uid + ":300000") //300000注销成功
	w.Write(rst)
}

func infoUserController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value
	if req.Method == "POST" {
		info := userinfo(uid)
		commentData, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			w.Write([]byte("600001"))
			return
		}
		io.Copy(w, bytes.NewReader(commentData))
	}
}

func addItemController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value
	if req.Method == "POST" {
		obj_name := req.FormValue("obj_name")
		obj_price := req.FormValue("obj_price")
		obj_info := req.FormValue("obj_info")
		use_time := req.FormValue("use_time")
		rst := addItem(obj_name, uid, obj_price, obj_info, use_time)
		w.Write(rst)
	}
}

func listItemController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "GET" {
		rst := listItem()
		commentData, err := json.MarshalIndent(rst, "", "    ")
		if err != nil {
			w.Write([]byte("600001"))
			return
		}
		io.Copy(w, bytes.NewReader(commentData))
	}
}

func shareRequestController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid_request := cookie_read.Value
	if req.Method == "POST" {
		uid_response := req.FormValue("uid_response")
		obj_name := req.FormValue("obj_name")
		obj_price := req.FormValue("obj_price")
		obj_info := req.FormValue("obj_info")
		use_time := req.FormValue("use_time")
		rst := shareRequest(uid_request, uid_response, obj_name, obj_price, obj_info, use_time)
		w.Write(rst)
	}
}

func shareResponseController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid_response := cookie_read.Value
	if req.Method == "POST" {
		uid_request := req.FormValue("uid_request")
		obj_name := req.FormValue("obj_name")
		obj_price := req.FormValue("obj_price")
		obj_info := req.FormValue("obj_info")
		use_time := req.FormValue("use_time")
		agree := req.FormValue("agree")
		rst := shareResponse(uid_request, uid_response, obj_name, obj_price, obj_info, use_time, agree)
		w.Write(rst)
	}
}

func sendMessageController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	from := cookie_read.Value
	if req.Method == "POST" {
		content := req.FormValue("content")
		to := req.FormValue("to")
		rst := sendMessage(content, from, to)
		w.Write(rst)
	}
}

func receiveMessageController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value
	if req.Method == "GET" {
		rst := receiveMessage(uid)
		commentData, err := json.MarshalIndent(rst, "", "    ")
		if err != nil {
			w.Write([]byte("600001"))
			return
		}
		io.Copy(w, bytes.NewReader(commentData))
	}
}

func listMessageController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value
	if req.Method == "GET" {
		rst := listMessage(uid)
		commentData, err := json.MarshalIndent(rst, "", "    ")
		if err != nil {
			w.Write([]byte("600001"))
			return
		}
		io.Copy(w, bytes.NewReader(commentData))
	}
}
