// Items.go
// Author : Faldict/cmc_iris
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

const (
	NOT_SUCCESS = iota
	SUCCESS
	UNDEFINE
)

func addItem(obj_name string, uid string, obj_price string, obj_info string, use_time string) []byte {
	upload_time := time.Now()
	obj_state := IN_TRADE

	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
		return []byte("300001") //300001OPEN错误
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Items VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return []byte("300002") //prepare错误
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(obj_name, uid, upload_time, string(obj_state), obj_price, obj_info, use_time) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("300003") //exec错误
	}

	return []byte("400000") //400000添加成功
}

func userinfo(uid string) []byte {
	db, err := sql.Open("mysql", "user:password@/datebase")
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
			return []byte("300005") //读取错误
		}
		rst = append(rst, tmp)
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("600001") // json错误
	}
	return b
}

func listItem() []byte {
	db, err := sql.Open("mysql", "user:password@/datebase")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Items")
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
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
		return []byte("600001")
	}
	return b
}

func shareRequest(uid_request string, uid_response string, obj_name string, obj_price string, obj_info string, use_time string) []byte {
	handle := "0"

	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
		return []byte("300001") //300001OPEN错误
	}

	stmtIns, err := db.Prepare("INSERT INTO ShareRequests VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return []byte("300002") //prepare错误
	}

	_, err = stmtIns.Exec(uid_request, uid_response, obj_name, obj_price, obj_info, use_time, handle, string(UNDEFINE)) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("300003") //exec错误
	}

	return []byte("400000") //400000添加成功
}

func shareResponse(uid_request string, uid_response string, obj_name string, obj_price string, obj_info string, use_time string, agree string) []byte {

	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}

	stmtUpd1, err := db.Prepare("UPDATE ShareRequests SET handle=? WHERE OBJ_name=? AND uid_request=? AND uid_response=?")
	if err != nil {
		log.Fatal(err)
		return []byte("300002")
	}

	_, err = stmtUpd1.Exec("1", obj_name, uid_request, uid_response)
	if err != nil {
		log.Fatal(err)
		return []byte("300003")
	}

	if agree == "1" {
		stmtUpd2, err := db.Prepare("UPDATE ShareRequests SET success=? WHERE OBJ_name=? AND uid_request=? AND uid_response=?")
		if err != nil {
			log.Fatal(err)
			return []byte("300002")
		}

		_, err = stmtUpd2.Exec(string(SUCCESS), obj_name, uid_request, uid_response)
		if err != nil {
			log.Fatal(err)
			return []byte("300003")
		}

		stmtUpd3, err := db.Prepare("UPDATE Items SET OBJ_state=? WHERE OBJ_name=? AND UID=?")
		if err != nil {
			log.Fatal(err)
			return []byte("300002")
		}

		_, err = stmtUpd3.Exec(string(NOT_IN_TRADE), obj_name, uid_response)
		if err != nil {
			log.Fatal(err)
			return []byte("300003")
		}
	} else {
		stmtUpd2, err := db.Prepare("UPDATE ShareRequests SET success=? WHERE OBJ_name=? AND uid_request=? AND uid_response=?")
		if err != nil {
			log.Fatal(err)
			return []byte("300002")
		}

		_, err = stmtUpd2.Exec(string(NOT_SUCCESS), obj_name, uid_request, uid_response)
		if err != nil {
			log.Fatal(err)
			return []byte("300003")
		}
	}

	return []byte("500000") //share成功
}
