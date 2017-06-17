// message.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func sendMessage(content string, from string, to string) []byte {
	upload_time := time.Now()
	read_state := "0"

	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001") //300001 OPEN错误
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Message VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return []byte("300002") //prepare错误
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(content, from, to, upload_time, read_state) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("300003") //exec错误
	}

	return []byte("400000") //400000添加成功
}

func receiveMessage(uid string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Message WHERE to=?", uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
	}

	type data struct {
		content     string
		from        string
		to          string
		upload_time string
	}

	var tmp data
	var read_state string
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.content, &tmp.from, &tmp.to, &tmp.upload_time, &read_state)
		if err != nil {
			log.Fatal(err)
			return []byte("300005") //读取错误
		}
		if read_state == "0" {
			rst = append(rst, tmp)
		}
	}
	b, err := json.Marshal(rst)
	if err != nil {
		return []byte("600001") // json错误
	}
	return b
}

func listMessage(uid string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Message WHERE to=? OR from=?", uid, uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300004")
	}

	type data struct {
		content     string
		from        string
		to          string
		upload_time string
		read_state  string
	}

	var tmp data
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.content, &tmp.from, &tmp.to, &tmp.upload_time, &tmp.read_state)
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
