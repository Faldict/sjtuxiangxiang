// Items.go
// Author : Faldict
// Date : 2017-05-02
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	_ = iota
	IN_TRADE
	NOT_IN_TRADE
	TIME_OUT
)

func addItem(obj_name string, uid string, obj_price string, obj_info string, use_time string) []byte {
	upload_time := time.Now()
	obj_state := IN_TRADE

	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Items VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = db.Exec(obj_name, uid, upload_time, string(obj_state), obj_price, obj_info, use_time) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("AddItem error")
	}

	return []byte("AddItem successfully")
}

func userinfo(uid string) []byte {
	db, err := sql.Open("mysql", "user:password@/datebase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM INFO_table WHERE user_ID LIKE %s", uid) // 不知道能不能这样用
	if err != nil {
		log.Fatal(err)
	}

	type data struct {
		uid                string
		photo              string // should be jpeg 喵喵喵？
		description        string
		Age                string
		RelationshipStatus string // {single, inlove}
		Jaccount           string
		score              string
	}

	var tmp data
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.uid, &tmp.photo, &tmp.description, &tmp.Age, &tmp.RelationshipStatus, &tmp.Jaccount, &tmp.score)
		if err != nil {
			log.Fatal(err)
		}
		rst = append(rst, tmp)
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("json error")
	}
	return b
}

func listItem() []byte {
	db, err := sql.Open("mysql", "user:password@/datebase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Items")
	if err != nil {
		log.Fatal(err)
	}

	type data struct {
		obj_name    string
		uid         string
		upload_time string
		obj_state   string // sorry to change it to be string
		obj_price   string // it need to be a string
		obj_info    string
		use_time    string
	}

	var tmp data
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.obj_name, &tmp.uid, &tmp.upload_time, &tmp.obj_state, &tmp.obj_price, &tmp.obj_info, &tmp.use_time)
		if err != nil {
			log.Fatal(err)
		}
		rst = append(rst, tmp)
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("json error")
	}
	return b
}

func shareItem(obj_name string, uid string, obj_price string, obj_info string, use_time string) []byte {
	obj_state := NOT_IN_TRADE

	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmtUpd, err := db.Prepare("UPDATE Items SET OBJ_state=? WHERE OBJ_name=? AND UID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtUpd.Close()

	_, err = db.Exec(string(obj_state), obj_name, uid)
	if err != nil {
		log.Fatal(err)
		return []byte("AddItem error")
	}

	return []byte("AddItem successfully")
}
