package routes

import (
	"net/http"
	"web_final/app"
	"web_final/controllers"
)

func PostRoutes(app *app.App) {

	postController := controllers.PostController{
		App: app,
	}

	http.HandleFunc("/register-handle", postController.Register)
	http.HandleFunc("/login-handle", postController.Login)
}
