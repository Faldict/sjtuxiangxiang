/*
   This is the user model for sjtuxiangxiang. And the functions are required by
   user controller.

   Author: Faldict
   Date: 2017-04-28

*/

package main

import (
	"database/sql"
	"log"
	"encoding/json"
	"fmt"
	//"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// TODO
// const (
//     DB_NAME = ""
//     DB_USER = ""
//     DB_PASSWD = ""
// )

func register(username string, passwd string, email string, description string, Age string, RelationStatus string, Jaccount string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO users VALUES( ?, ?, ? )") // ? = placeholder
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(username, passwd, email) // Insert tuples (i, i^2)
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
		return []byte("Register Error")
	}

	stmtIns2, err := db.Prepare("INSERT INTO INFO_table VALUES( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmtIns2.Close()

	var score, num string
	score = "50"
	num = "0"
	_, err = stmtIns2.Exec(username, description, Age, RelationStatus, Jaccount, score, num) // Insert tuples (i, i^2)
	if err != nil {
		log.Fatal(err.Error())
		return []byte("Register Error")
	}

	return []byte("Register successfully")
}

func updateInfo(username string, description string, age string, RelationStatus string, Jaccount string) string {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("UPDATE INFO_table SET description = ?, age = ?, RelationshipStatus = ?, Jaccount = ? WHERE user_ID = ?")
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(description, age, RelationStatus, Jaccount, username) // Insert tuples (i, i^2)
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
		return "Update Error"
	}
	return "Update successfully!"
}

//func login(w http.ResponseWriter, uid string, passwd string) []byte {
/*
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	var PSD string
	err = db.QueryRow("SELECT user_PSD FROM users WHERE user_ID=?", uid).Scan(&PSD)
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300002")
	}

	if passwd == PSD {
		cookie := http.Cookie{
			Name:   "uid",
			Value:  uid,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(w, &cookie)
		return []byte("200000")
	} else {
		return []byte("300001")
	}
*/
//}

//func logout(w http.ResponseWriter, req http.Request) []byte {
/*
	uid, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		return []byte("Cookie read error")
	}
*/
//var uid string
//for _, cookie := range req.Cookies() {
//	uid = cookie.Name
//	break
//}

//cookie := http.Cookie{
//	Name:   "uid",
//	Value:  uid,
//	Path:   "/",
//	MaxAge: -1,
//}
//http.SetCookie(w, &cookie)

//return []byte("Logout successfully")
//}

func tradeRecord(uid string) []byte {
	fmt.Println("uid : " + uid)
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ShareRequests WHERE uid_request = ? OR uid_response = ?", uid, uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
	}

	type data struct {
		Obj_name   string // 对方的uid
		Uid_other    string
		Cnt         string
		Upload_time string
		Typ         string // "0"/"1" 需求出租
	}

	var tmp data
	var uid_1, uid_2 string
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.Obj_name, &uid_1, &uid_2, &tmp.Cnt, &tmp.Upload_time)
		if err != nil {
			log.Fatal(err)
		}
		if uid_1 == uid {
			tmp.Uid_other = uid_2
			tmp.Typ = "0"
			fmt.Println("uid1 : " + uid_1)
		} else {
			tmp.Uid_other = uid_1
			tmp.Typ = "1"
			fmt.Println("uid2 : " + uid_2)
		}
		rst = append(rst, tmp)
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("600001")
	}
	return b
}

func listShare(uid string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Items WHERE Uploader = ?", uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
	}

	type data struct {
		Obj_name    string
		Uid         string
		Upload_time string
		Obj_state   string // sorry to change it to be string
		Obj_price   string // it need to be a string
		Obj_info    string
		Use_time    string
		Obj_type    string
	}

	var tmp data
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.Obj_name, &tmp.Uid, &tmp.Upload_time, &tmp.Obj_state, &tmp.Obj_price, &tmp.Obj_info, &tmp.Use_time, &tmp.Obj_type)
		if err != nil {
			log.Fatal(err)
		}
		rst = append(rst, tmp)
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("600001")
	}
	return b
}

func userInfo(uid string) []byte {
	fmt.Println(uid)
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM INFO_table WHERE user_ID=?", uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
	}

	type data struct {
		Uid string
		// photo
		Description        string
		Age                string
		RelationshipStatus string // {single, inlove}
		Jaccount           string
		Score              string
		Num                string
		Phone              string
	}

	var tmp data

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.Uid, &tmp.Description, &tmp.Age, &tmp.RelationshipStatus, &tmp.Jaccount, &tmp.Score, &tmp.Num)
		if err != nil {
			log.Fatal(err)
			return []byte("300005") //读取错误
		}
	}

	db.QueryRow("SELECT EMAIL FROM users WHERE user_ID = ?", uid).Scan(&tmp.Phone)
	b, err := json.Marshal(tmp)
	if err != nil {
		return []byte("600001") // json错误
	}
	return b
}