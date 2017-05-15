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
	//"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// TODO
// const (
//     DB_NAME = ""
//     DB_USER = ""
//     DB_PASSWD = ""
// )

func register(username string, passwd string, email string) []byte {
	// TODO
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// TODO
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

	return []byte("Register successfully")
}

//func login(w http.ResponseWriter, uid string, passwd string) []byte {
/*
	db, err := sql.Open("mysql", "user:password@/dbname")
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
