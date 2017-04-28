/*
    This is the user model for sjtuxiangxiang. And the functions are required by
    user controller.

    Author: Faldict
    Date: 2017-04-28

*/

package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

// TODO
// const (
//     DB_NAME = ""
//     DB_USER = ""
//     DB_PASSWD = ""
// )

func register(username string, passwd string, email string) {
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
    }
}
