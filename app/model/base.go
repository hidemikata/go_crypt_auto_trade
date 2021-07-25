package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
const format1 = "2006-01-02 15:04:05" //固定

func init() {
	dbcon, err := sql.Open("mysql", "root:@/coin_data")
	if err != nil {
		fmt.Println("db connect")
	}
	db = dbcon
	//    defer dbcon.Close()
}
