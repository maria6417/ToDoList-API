package database

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB
var err error

func Connect() {
	DB, err = sql.Open("mysql", "mhirai:WelCome##1234@tcp(18.118.134.59:3306)/test01")
	if err != nil {
		panic("could not connect to database")
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("failed connecting to database")
		return
	} else {
		fmt.Println("database connected!")
		return
	}

}
