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
	"net/http"

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

func login(w http.ResponseWriter, uid string, passwd string) []byte {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE user_ID LIKE %s", uid)
	if err != nil {
		log.Fatal(err.Error())
	}

	type data struct {
		user_ID   string
		user_PSD  string
		user_INFO string
	}

	var tmp data
	err = rows.Scan(&tmp.user_ID, &tmp.user_PSD, &tmp.user_INFO)
	if err != nil {
		log.Fatal(err)
	}

	if passwd == tmp.user_PSD {
		cookie := http.Cookie{
			Name:   "uid",
			Value:  uid,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(w, &cookie)

		db, err := sql.Open("mysql", "user:password@/dbname")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer db.Close() // defer的机制我不是很清楚 所以不知道这里是不是还要open一次

		stmtUpd, err := db.Prepare("UPDATE users SET user_INFO=? WHERE user_ID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmtUpd.Close()

		_, err = db.Exec("1", uid) // "1" means true, set "0" when initialize
		if err != nil {
			log.Fatal(err)
		}

		return []byte("Login successfully")
	} else {
		return []byte("PSD wrong")
	}
}

func logout(w http.ResponseWriter, req http.Request) []byte {
	uid, err := req.Cookie("uid")
	if err != nil {
		log.Fatal(err.Error())
		return []byte("Cookie read error")
	}

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	stmtUpd, err := db.Prepare("UPDATE users SET user_INFO=? WHERE user_ID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtUpd.Close()

	_, err = db.Exec("0", uid)
	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{
		Name:   "uid",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	return []byte("Logout successfully")
}
