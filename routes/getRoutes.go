package routes

import (
	"io/fs"
	"log"
	"net/http"
	"web_final/app"
	"web_final/controllers"
)

func GetRoutes(app *app.App) {

	getController := controllers.GetController{
		App: app,
	}

	staticFS, err := fs.Sub(app.Res, "static")
	if err != nil {
		log.Println(err)
		return
	}
	staticHandler := http.FileServer(http.FS(staticFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	log.Println("Archivos cargados correctamente.")

	http.HandleFunc("/", getController.ShowHome)
	http.HandleFunc("/login", getController.ShowLogin)
	http.HandleFunc("/register", getController.ShowRegister)
	http.HandleFunc("/logout", getController.Logout)
}
