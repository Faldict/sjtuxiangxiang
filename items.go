// Items.go
// Author : Faldict/cmc_iris
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
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
		uid string
		// photo
		description        string
		Age                string
		RelationshipStatus string // {single, inlove}
		Jaccount           string
		score              string
		num                string
	}

	var tmp data
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.uid, &tmp.description, &tmp.Age, &tmp.RelationshipStatus, &tmp.Jaccount, &tmp.score, &tmp.num)
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

func updateScore(obj_uid string, obj_score string) []byte {
	var cur_score, cur_num, new_score, new_num string
	var cscore, cnum, oscore, nscore, nnum int
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300001")
	}
	defer db.Close()

	err = db.QueryRow("SELECT score FROM INFO_table WHERE user_ID=?", obj_uid).Scan(&cur_score)
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300004") //300002SELECT错误
	} //取当前得分

	err = db.QueryRow("SELECT num FROM INFO_table WHERE user_ID=?", obj_uid).Scan(&cur_num)
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300004") //300002SELECT错误
	} //取当前评价人数

	cscore, err = strconv.Atoi(cur_score)
	if err != nil {
		log.Fatal(err)
		return []byte("字符串转换成整数失败")
	}
	cnum, err = strconv.Atoi(cur_num)
	if err != nil {
		log.Fatal(err)
		return []byte("字符串转换成整数失败")
	}
	oscore, err = strconv.Atoi(obj_score)
	if err != nil {
		log.Fatal(err)
		return []byte("字符串转换成整数失败")
	}
	nnum = cnum + 1
	nscore = (cscore*cnum + oscore) / (cnum + 1) //计算新得分
	new_score = strconv.Itoa(nscore)
	new_num = strconv.Itoa(nnum)

	stmtUpd1, err := db.Prepare("UPDATE INFO_table SET score=? WHERE user_ID=?")
	if err != nil {
		log.Fatal(err)
		return []byte("300002")
	}

	_, err = stmtUpd1.Exec(new_score, obj_uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300003")
	} //更新得分

	stmtUpd2, err := db.Prepare("UPDATE INFO_table SET num=? WHERE user_ID=?")
	if err != nil {
		log.Fatal(err)
		return []byte("300002")
	}

	_, err = stmtUpd2.Exec(new_num, obj_uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300003")
	} //评价人数+1

	return []byte("评价成功")
}
