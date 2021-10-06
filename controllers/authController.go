package controller

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/easyTodoList/database"
	"example.com/easyTodoList/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("using-my-secret-key")

func SignUp(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := hashPassword(r.FormValue("password"))

	fmt.Printf("username: %s, password: %s", username, password)
	doesUserExist := checkUser(username)
	if doesUserExist {
		// should i be writing to response headers??
		io.WriteString(w, "username already taken.")
		return
	}

	stmt, err := database.DB.Prepare("insert into users values (?, ? )")
	checkErr(err)

	// Q: what am i supposed to do with results?
	results, err := stmt.Exec(username, string(password))
	checkErr(err)

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 1 {
		w.Header().Set("Content-Header", "application/json")
		io.WriteString(w, `{"message": "user registered"}`)
		return
	} else if rowsAffected == 0 {
		w.Header().Set("Content-Header", "application/json")
		io.WriteString(w, `{"message": "user failed to register"}`)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	doesUserExist := checkUser(username)
	if !doesUserExist {
		w.Header().Set("Content-Header", "application/json")
		io.WriteString(w, `{"message": "user does not exist"}`)
		return
	}

	err := checkPassword(username, password)
	if err != nil {
		w.Header().Set("Content-Header", "application/json")
		io.WriteString(w, `{"message": "password does not match"}`)
		return
	}
	fmt.Println("user validated.")

	//register token on cookie.
	// first generate token.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString(secretKey)
	if err != nil {
		w.Header().Set("Content-Header", "application/json")
		io.WriteString(w, `{"message": "could not login"}`)
		checkErr(err)
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Header", "application/json")
	io.WriteString(w, `{"message": "logged in successfully!"}`)

}

func hashPassword(pass string) []byte {
	password, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	checkErr(err)

	return password
}

func checkUser(username string) bool {
	rows, err := database.DB.Query("select count(*) from users where username = ? ", username)
	checkErr(err)

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}

	if count == 0 {
		fmt.Println("user does not exist")
		return false
	} else {
		fmt.Println("user exists")
		return true
	}
}

func checkPassword(username, password string) error {
	user := getUserData(username)
	fmt.Println("password passed : ", password)
	fmt.Println("password in database:", user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}

func getUserData(username string) models.User {

	var user models.User
	fmt.Println("getting user data for username:", username)
	rows, err := database.DB.Query("select * from users where username = ?", username)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&user.Username, &user.Password)
		checkErr(err)
		fmt.Println("user retrieved. username:", user.Username)
	}

	return user

}

// returns username if authenticated
func AuthenticateToken(r *http.Request) (string, error) {
	var user string
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return user, err
	}

	// get cookie value
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if token.Valid {
		fmt.Println("token is valid")

		// get user data from database
		claims := token.Claims.(jwt.StandardClaims)
		username := claims.Issuer
		return username, err

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return user, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Token expired.")
			return user, err
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return user, err
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return user, err
	}

}
