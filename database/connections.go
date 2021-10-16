package database

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB
var err error

var id int

func Connect() {
	DB, err = sql.Open("mysql", "mhirai:Test1234##@tcp(localhost:3306)/mhirai")
	if err != nil {
		panic("could not connect to database")
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("failed connecting to database")
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("database connected!")
		return
	}

}
