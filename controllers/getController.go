package controllers

import (
	"log"
	"net/http"
	"web_final/app"
	"web_final/models"
	"web_final/security"
	"web_final/templating"
)

// GetController is a wrapper struct for the App struct
type GetController struct {
	App *app.App
}

func (getController *GetController) ShowHome(w http.ResponseWriter, _ *http.Request) {
	type dataStruct struct {
		Test string
	}

	data := dataStruct{
		Test: "Hola mundo mundial!",
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/home.html", data)
}

func (getController *GetController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	// Create csrf token
	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		log.Println("Error al generar el token.")
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/register.html", data)
}

func (getController *GetController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	// Create csrf token
	CsrfToken, err := security.GenerateCsrfToken(w, r)
	if err != nil {
		log.Println("Error al generar el token.")
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/login.html", data)
}

func (getController *GetController) Logout(w http.ResponseWriter, r *http.Request) {
	models.LogoutUser(getController.App, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
