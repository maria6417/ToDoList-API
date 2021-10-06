package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"example.com/easyTodoList/database"
	"example.com/easyTodoList/models"
)

func GetTodoItem(w http.ResponseWriter, r *http.Request) {

	username, err := AuthenticateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	var data []byte
	id := r.URL.Query().Get("id")
	if id != "" {
		i, err := strconv.Atoi(id)
		checkErr(err)
		// get by Id
		toDo, err := GetItemById(i, username)
		checkErr(err)
		data, err = json.MarshalIndent(toDo, "", " ")
		checkErr(err)
	} else {
		toDo, err := GetAllItems(username)
		checkErr(err)
		data, err = json.MarshalIndent(toDo, "", " ")
		checkErr(err)
	}

	io.WriteString(w, string(data))
}

func GetItemById(id int, username string) (models.ToDo, error) {
	var todo models.ToDo
	rows, err := database.DB.Query("select id, title, description, completed from todo where id = ? and username = ?", id, username)
	checkErr(err)

	err = errors.New("no data found")
	for rows.Next() {
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
	}
	return todo, err
}

func GetAllItems(username string) ([]models.ToDo, error) {

	var toDos []models.ToDo
	var err error

	rows, err := database.DB.Query("SELECT id, title, description, completed FROM todo where username = ?", username)
	if err != nil {
		return toDos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.ToDo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return toDos, err
		}
		toDos = append(toDos, todo)
	}

	return toDos, err
}

func GetItemsByStatus(status bool, username string) ([]models.ToDo, error) {

	var toDos []models.ToDo
	var err error

	rows, err := database.DB.Query("SELECT id, title, description, completed FROM todo where completed = ? and username = ?", status, username)
	if err != nil {
		return toDos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.ToDo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return toDos, err
		}
		toDos = append(toDos, todo)
	}

	return toDos, err
}

func GetNewId(username string) int {

	var id int
	rows, err := database.DB.Query("select max(id)+1 from todo where username = ?", username)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		checkErr(err)
		break
	}
	return id
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {

	username, err := AuthenticateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	// get all data where completed column is true
	toDos, err := GetItemsByStatus(true, username)
	checkErr(err)

	data, err := json.MarshalIndent(toDos, "", " ")
	if err != nil {
		panic(err.Error())
	}
	io.WriteString(w, string(data))

}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {

	username, err := AuthenticateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	// get all data where completed column is false
	toDos, err := GetItemsByStatus(false, username)
	checkErr(err)

	data, err := json.MarshalIndent(toDos, "", " ")
	if err != nil {
		panic(err.Error())
	}
	io.WriteString(w, string(data))
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

	username, err := AuthenticateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

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

	_, err = GetItemById(id, username)
	checkErr(err)

	stmt, err := database.DB.Prepare(`update todo set completed = ? where id = ? and username = ?`)
	checkErr(err)
	defer stmt.Close()

	completed := r.FormValue("completed")

	fmt.Println("id ", id)
	fmt.Println("completed", completed)

	rs, err := stmt.Exec(completed == "true", id, username)
	checkErr(err)

	ra, err := rs.RowsAffected()
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated" : true , "updatedRows" : `+strconv.FormatInt(ra, 10)+`}`)

}

func InsertItem(w http.ResponseWriter, r *http.Request) {

	username, err := AuthenticateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	id := GetNewId(username)
	title := r.FormValue("title")
	description := r.FormValue("description")
	completed := r.FormValue("completed") == "true"

	stmt, err := database.DB.Prepare(`insert into todo (id, title, description, completed, username ) values (?, ?, ?, ? ,?)`)
	checkErr(err)
	defer stmt.Close()

	rs, err := stmt.Exec(id, title, description, completed, username)
	checkErr(err)

	ra, err := rs.RowsAffected()
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"inserted" : true , "insertedRows" : `+strconv.FormatInt(ra, 10)+`}`)

}
