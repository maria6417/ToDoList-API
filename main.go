package main

import (
	"net/http"

	"github.com/gorilla/mux"

	database "example.com/easyTodoList/database"
	route "example.com/easyTodoList/routes"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	database.Connect()

	mux := mux.NewRouter()
	route.SetRoutes(mux)
	http.ListenAndServe(":8080", mux)

}
