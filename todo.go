package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type ToDo struct {
	Id          int
	Title       string
	Description string
	Completed   bool
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "mhirai:WelCome##1234@tcp(18.118.134.59:3306)/test01")
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()

	if err != nil {
		fmt.Println("データベース接続失敗")
		return
	} else {
		fmt.Println("データベース接続成功")
	}
}

func main() {
	http.HandleFunc("/api/getTodoItem", getTodoItem)
	http.HandleFunc("/api/getCompletedItems", getCompletedItems)
	http.HandleFunc("/api/getIncompleteItems", getIncompleteItems)
	http.HandleFunc("/api/updateItem", updateItem)
	http.HandleFunc("/api/insertItem", insertItem)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func getTodoItem(w http.ResponseWriter, r *http.Request) {

	var data []byte
	id := r.URL.Query().Get("id")
	if id != "" {
		i, err := strconv.Atoi(id)
		checkErr(err)
		// get by Id
		toDo, err := getItemById(i)
		checkErr(err)
		data, err = json.MarshalIndent(toDo, "", " ")
		checkErr(err)
	} else {
		toDo, err := getAllItems()
		checkErr(err)
		data, err = json.MarshalIndent(toDo, "", " ")
		checkErr(err)
	}

	io.WriteString(w, string(data))
}

func getCompletedItems(w http.ResponseWriter, r *http.Request) {

	// get all data where completed column is true
	toDos, err := getItemsByStatus(true)
	checkErr(err)

	data, err := json.MarshalIndent(toDos, "", " ")
	if err != nil {
		panic(err.Error())
	}
	io.WriteString(w, string(data))

}

func getIncompleteItems(w http.ResponseWriter, r *http.Request) {
	var toDos []ToDo
	// get all data where completed column is false
	toDos, err := getItemsByStatus(false)
	checkErr(err)

	data, err := json.MarshalIndent(toDos, "", " ")
	if err != nil {
		panic(err.Error())
	}
	io.WriteString(w, string(data))
}

func updateItem(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		err := errors.New("method is not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	// /api/updateItem?id=1 --data
	// check if item id exists in db
	u := r.URL
	qs := u.Query()
	id, err := strconv.Atoi(qs.Get("id"))
	checkErr(err)

	_, err = getItemById(id)
	checkErr(err)

	stmt, err := db.Prepare(`update todo set completed = ? where id = ?`)
	checkErr(err)
	defer stmt.Close()

	completed := r.FormValue("completed")

	fmt.Println("id ", id)
	fmt.Println("completed", completed)

	rs, err := stmt.Exec(completed == "true", id)
	checkErr(err)

	ra, err := rs.RowsAffected()
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated" : true , "updatedRows" : `+strconv.FormatInt(ra, 10)+`}`)

}

func insertItem(w http.ResponseWriter, r *http.Request) {

	id := getNewId()
	title := r.FormValue("title")
	description := r.FormValue("description")
	completed := r.FormValue("completed") == "true"

	stmt, err := db.Prepare(`insert into todo values (?, ?, ? , ?)`)
	checkErr(err)
	defer stmt.Close()

	rs, err := stmt.Exec(id, title, description, completed)
	checkErr(err)

	ra, err := rs.RowsAffected()
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"inserted" : true , "insertedRows" : `+strconv.FormatInt(ra, 10)+`}`)

}
