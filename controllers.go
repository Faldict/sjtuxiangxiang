package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
	case "tradeRecord":
		tradeRecordController(w, req)
	case "listShare":
		listShareController(w, req)
	case "updateinfo":
		updateInfoController(w, req)
	case "information":
		informationController(w, req)
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
	case "updateScore":
		updateScoreController(w, req)
	case "info":
		infoController(w, req)
	default:
		NotFound(w, req)
	}
}

func registerUserController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		username := req.FormValue("username")
		passwd := req.FormValue("password")
		email := req.FormValue("email")
		//need more userInfo
		description := req.FormValue("description")
		Age := req.FormValue("Age")
		RelationStatus := req.FormValue("RelationStatus")
		Jaccount := req.FormValue("Jaccount")
		rst := register(username, passwd, email, description, Age, RelationStatus, Jaccount)
		w.Write(rst)
	}
}

func updateInfoController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		cookie_read, err := req.Cookie("uid")
		if err != nil {
			log.Fatal(err)
			io.WriteString(w, "Not login!")
			return 
		}

		username := cookie_read.Value
		description := req.FormValue("description")
		Jaccount := req.FormValue("jaccount")
		age := req.FormValue("age")
		RelationStatus := req.FormValue("RelationStatus")
		result := updateInfo(username, description, age, RelationStatus, Jaccount)
		io.WriteString(w, result)
	}
}

func loginUserController(w http.ResponseWriter, req *http.Request) {
	var rst []byte
	var PSD string
	if req.Method == "POST" {
		username := req.FormValue("username")
		passwd := req.FormValue("password")
		db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
		if err != nil {
			log.Fatal(err.Error())
			rst = []byte("300001")
			w.Write(rst)
			return
		}
		defer db.Close()

		err = db.QueryRow("SELECT user_PSD FROM users WHERE user_ID=?", username).Scan(&PSD)
		if err != nil {
			log.Fatal(err.Error())
			rst = []byte("300004") //300002SELECT错误
			w.Write(rst)
			return
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
			w.Write(rst)
			fmt.Println("login success!")
			return
		} else {
			rst = []byte("200001") //300001密码错误
			w.Write(rst)
			return
		}
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
	if req.Method == "GET" {
		info := userInfo(uid)
		w.Write(info)
	}
}

func informationController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == "POST" {
		userid := req.FormValue("id")
		info := userInfo(userid)
		w.Write(info)
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
		use_time := req.FormValue("end_time")
		typ := req.FormValue("type")
		rst := addItem(obj_name, uid, obj_price, obj_info, use_time, typ)
		w.Write(rst)
	}
}

func listItemController(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "GET" {
		query := req.URL.Query()
		typ := query["type"][0]
		rst := listItem(typ)
		w.Write(rst)
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
		rst := shareRequest(uid_request, uid_response, obj_name)
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
		agree := req.FormValue("agree") //to be "1" or "0"
		rst := shareResponse(uid_request, uid_response, obj_name, agree)
		w.Write(rst)
	}
}

func infoController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		obj_name := req.FormValue("id")
		rst := itemInfo(obj_name)
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

func updateScoreController(w http.ResponseWriter, req *http.Request) {
	cookie_read, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		w.Write([]byte("100000")) //100000未登录
		return
	}
	uid := cookie_read.Value //要登陆才能评价
	fmt.Println(uid)
	if req.Method == "POST" {
		obj_uid := req.FormValue("obj_uid")
		obj_score := req.FormValue("obj_score")
		rst := updateScore(obj_uid, obj_score)
		w.Write(rst)
	}
}

func tradeRecordController(w http.ResponseWriter, req *http.Request) {
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
		rst := tradeRecord(uid)
		w.Write(rst)
	}
}

func listShareController(w http.ResponseWriter, req *http.Request) {
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
		rst := listShare(uid)
		w.Write(rst)
	}
}
