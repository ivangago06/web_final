package routes

import (
	"net/http"
	"web_final/app"
	"web_final/controllers"
)

// PostRoutes defines all project post routes
func PostRoutes(app *app.App) {
	// Post controller struct initialize
	postController := controllers.PostController{
		App: app,
	}
	// User authentication
	http.HandleFunc("/register-handle", postController.Register)
	http.HandleFunc("/login-handle", postController.Login)
}
