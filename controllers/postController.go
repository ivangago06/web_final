package controllers

import (
	"log"
	"net/http"
	"time"
	"web_final/app"
	"web_final/models"
	"web_final/security"
)

// PostController is a wrapper struct for the App struct
type PostController struct {
	App *app.App
}

func (postController *PostController) Login(w http.ResponseWriter, r *http.Request) {
	// Validate csrf token
	_, err := security.VerifyCsrfToken(r)
	if err != nil {
		log.Println("Error verificando csrf token")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	if username == "" || password == "" {
		log.Println("Intente iniciar sesión con un nombre de usuario o una contraseña.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	_, err = models.AuthenticateUser(postController.App, w, username, password, remember)
	if err != nil {
		log.Println("Error de autenticación de usuario")
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (postController *PostController) Register(w http.ResponseWriter, r *http.Request) {
	// Validate csrf token
	_, err := security.VerifyCsrfToken(r)
	if err != nil {
		log.Println("Error verificando csrf token")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	createdAt := time.Now()
	updatedAt := time.Now()

	if username == "" || password == "" {
		log.Println("Intente iniciar sesión con un nombre de usuario o una contraseña.")
		http.Redirect(w, r, "/register", http.StatusFound)
	}

	_, err = models.CreateUser(postController.App, username, password, createdAt, updatedAt)
	if err != nil {
		log.Println("Error al crear el usuario")
		log.Println(err)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
