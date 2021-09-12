package main

import (
	"database/sql"
	"errors"
)

func init() {
	db, err = sql.Open("mysql", "mhirai:WelCome##1234@tcp(18.118.134.59:3306)/test01")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getItemById(id int) (ToDo, error) {
	var err error
	var todo ToDo
	rows, err := db.Query("select * from todo where id = ?", id)
	checkErr(err)

	err = errors.New("no data found")
	for rows.Next() {
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
	}
	return todo, err
}

func getAllItems() ([]ToDo, error) {

	var toDos []ToDo
	var err error

	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		return toDos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo ToDo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return toDos, err
		}
		toDos = append(toDos, todo)
	}

	return toDos, err
}

func getItemsByStatus(status bool) ([]ToDo, error) {

	var toDos []ToDo
	var err error

	rows, err := db.Query("SELECT * FROM todo where completed = ?", status)
	if err != nil {
		return toDos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo ToDo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return toDos, err
		}
		toDos = append(toDos, todo)
	}

	return toDos, err
}

func getNewId() int {

	var id int
	rows, err := db.Query("select max(id)+1 from todo")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		checkErr(err)
		break
	}
	return id
}
