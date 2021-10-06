package models

type ToDo struct {
	Id          int
	Title       string
	Description string
	Completed   bool
}

type User struct {
	Username string
	Password string
}
