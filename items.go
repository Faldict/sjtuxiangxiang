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

func addItem(obj_name string, uid string, obj_price string, obj_info string, use_time string, typ string) []byte {
	upload_time := time.Now()
	obj_state := IN_TRADE

	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001") //300001OPEN错误
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Items VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return []byte("300002") //prepare错误
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(obj_name, uid, upload_time, string(obj_state), obj_price, obj_info, use_time, typ) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("300003") //exec错误
	}

	return []byte("400000") //400000添加成功
}

func listItem(typ string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Items WHERE obj_type = ?", typ)
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
	}

	var tmp data
	var tmp_type string
	rst := []data{}

	for rows.Next() {
		rows.Columns()
		err = rows.Scan(&tmp.Obj_name, &tmp.Uid, &tmp.Upload_time, &tmp.Obj_state, &tmp.Obj_price, &tmp.Obj_info, &tmp.Use_time, &tmp_type)
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

func shareRequest(uid_request string, uid_response string, obj_name string) []byte {
	upload_time := time.Now()

	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001") //300001OPEN错误
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO ShareRequests VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return []byte("300002") //prepare错误
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(obj_name, uid_request, uid_response, "0", upload_time) // obj_state is string
	if err != nil {
		log.Fatal(err)
		return []byte("300003") //exec错误
	}

	return []byte("400000") //400000添加成功
}

func shareResponse(uid_request string, uid_response string, obj_name string, agree string) []byte {

	var count string
	var cnt int

	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err)
		return []byte("300001")
	}
	defer db.Close()

	err = db.QueryRow("SELECT cnt FROM ShareRequests WHERE OBJ_name=? AND uid_request=? AND uid_response=?", obj_name, uid_request, uid_response).Scan(&count)
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300004") //300002SELECT错误
	}

	cnt, err = strconv.Atoi(count)
	if err != nil {
		log.Fatal(err)
		return []byte("500001")
	}

	if agree == "1" {
		cnt = cnt + 1
		count = strconv.Itoa(cnt)
		stmtUpd1, err := db.Prepare("UPDATE ShareRequests SET cnt=? WHERE OBJ_name=? AND uid_request=? AND uid_response=?")
		if err != nil {
			log.Fatal(err)
			return []byte("300002")
		}
		defer stmtUpd1.Close()

		_, err = stmtUpd1.Exec(count, obj_name, uid_request, uid_response)
		if err != nil {
			log.Fatal(err)
			return []byte("300003")
		}
	}
	if cnt == 2 {
		stmtUpd2, err := db.Prepare("UPDATE Items SET OBJ_state=? WHERE OBJ_name=? AND UID=?")
		if err != nil {
			log.Fatal(err)
			return []byte("300002")
		}
		defer stmtUpd2.Close()

		_, err = stmtUpd2.Exec(string(NOT_IN_TRADE), obj_name, uid_response)
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
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
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
		return []byte("500001")
	}
	cnum, err = strconv.Atoi(cur_num)
	if err != nil {
		log.Fatal(err)
		return []byte("500001")
	}
	oscore, err = strconv.Atoi(obj_score)
	if err != nil {
		log.Fatal(err)
		return []byte("500001")
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
	defer stmtUpd1.Close()

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
	defer stmtUpd2.Close()

	_, err = stmtUpd2.Exec(new_num, obj_uid)
	if err != nil {
		log.Fatal(err)
		return []byte("300003")
	} //评价人数+1

	return []byte("400000")
}

func itemInfo(obj_id string) []byte {
	db, err := sql.Open("mysql", "sjtuxx:sjtuxx@tcp(localhost:3306)/sjtuxiangxiang")
	if err != nil {
		log.Fatal(err.Error())
		return []byte("300001")
	}
	defer db.Close()

	type data struct {
		Obj_name    string
		Uploader    string
		Upload_time string
		Obj_price   string // it need to be a string
		Obj_info    string
		Use_time    string
		Score       int
	}
	var result data

	err = db.QueryRow("SELECT Uploader, UploadTime, OBJ_price, OBJ_INFO, OBJ_usetime FROM Items WHERE OBJ_name = '" + obj_id + "'").Scan(&result.Uploader, &result.Upload_time, &result.Obj_price, &result.Obj_info, &result.Use_time)
	result.Score = 80
	// err = db.QueryRow("SELECT score FROM info WHERE user_ID = ?", result.uploader).Scan(&result.score)
	if err != nil {
		log.Fatal(err)
		return []byte("300002")
	}

	result.Obj_name = obj_id
	j, err := json.Marshal(result)
	return j
}
