package models

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"web_final/app"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const userColumnsNoId = "\"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\""
const userColumns = "\"Id\", " + userColumnsNoId
const userTable = "public.\"User\""

const (
	selectUserById       = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Id\" = $1"
	selectUserByUsername = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Username\" = $1"
	insertUser           = "INSERT INTO " + userTable + " (" + userColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
)

func GetCurrentUser(app *app.App, r *http.Request) (User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error al obtener la cookie de sesión")
		return User{}, err
	}

	session, err := GetSessionByAuthToken(app, cookie.Value)
	if err != nil {
		log.Println("Error al obtener sesión por auth token")
		return User{}, err
	}

	return GetUserById(app, session.UserId)
}

func GetUserById(app *app.App, id int64) (User, error) {
	user := User{}

	// Query row by id
	err := app.Db.QueryRow(selectUserById, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Usuario no encontrado para el id:" + strconv.FormatInt(id, 10))
		return User{}, err
	}

	return user, nil
}

func GetUserByUsername(app *app.App, username string) (User, error) {
	user := User{}

	// Query row by username
	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Usuario no encontrado para el usuario:" + username)
		return User{}, err
	}

	return user, nil
}

func CreateUser(app *app.App, username string, password string, createdAt time.Time, updatedAt time.Time) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error de hashing de contraseña al crear usuario")
		return User{}, err
	}

	var lastInsertId int64

	err = app.Db.QueryRow(insertUser, username, string(hash), createdAt, updatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error al crear el registro de usuario")
		return User{}, err
	}

	return GetUserById(app, lastInsertId)
}

func AuthenticateUser(app *app.App, w http.ResponseWriter, username string, password string, remember bool) (Session, error) {
	var user User

	// Query row by username
	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error de autenticación para la usuario:" + username)
		return Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Error de autenticación, password incorrecto para la usuario:" + username)
		return Session{}, err
	} else {
		return CreateSession(app, w, user.Id, remember)
	}
}

func LogoutUser(app *app.App, w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error al obtener la cookie de la solicitud")
		return
	}

	err = DeleteSessionByAuthToken(app, w, cookie.Value)
	if err != nil {
		log.Println("Error al eliminar sesión por AuthToken")
		return
	}
}
