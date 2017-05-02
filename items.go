// Items.go
// Author : Faldict
// Date : 2017-05-02
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "time"
    "encoding/json"
)

const (
    _ = iota
    IN_TRADE
    NOT_IN_TRADE
    TIME_OUT
)

func addItem(obj_name string, uid string, obj_price int, obj_info string, use_time string) int {
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

    _, err := db.Exec(obj_name, uid, upload_time, obj_state, obj_price, obj_info, use_time)
    if err != nil {
        log.Fatal(err)
        return 0
    }

    return 1
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
        obj_name string
        uid string
        upload_time string
        obj_state int
        obj_price int
        obj_info string
        use_time string
    }

    var tmp data
    rst := []data {}

    for rows.Next() {
        rows.Columns()
        err = rows.Scan(&tmp.obj_name, &tmp.uid, &tmp.upload_time, &tmp.obj_state, &tmp.obj_price, &tmp.obj_info, &tmp.use_time)
        if err != nil {
            log.Fatal(err)
        }
        rst := append(rst, tmp)
    }
    return json.Marshal(rst)
}
