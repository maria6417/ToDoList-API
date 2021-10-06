package route

import (
	"net/http"

	controller "example.com/easyTodoList/controllers"
	"github.com/gorilla/mux"
)

func SetRoutes(mux *mux.Router) {
	mux.HandleFunc("/api/signup", controller.SignUp)
	mux.HandleFunc("/api/login", controller.Login)
	mux.HandleFunc("/api/getTodoItem", controller.GetTodoItem)
	mux.HandleFunc("/api/getCompletedItems", controller.GetCompletedItems)
	mux.HandleFunc("/api/getIncompleteItems", controller.GetIncompleteItems)
	mux.HandleFunc("/api/updateItem", controller.UpdateItem)
	mux.HandleFunc("/api/insertItem", controller.InsertItem)
	mux.Handle("/favicon.ico", http.NotFoundHandler())
}
