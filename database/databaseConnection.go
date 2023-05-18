package database

import (
	"database/sql"
	"fmt"
	"log"
	"web_final/app"

	_ "github.com/lib/pq"
)

func ConnectDB(app *app.App) *sql.DB {

	postgresConfig := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		app.Config.Db.Ip, app.Config.Db.Port, app.Config.Db.User, app.Config.Db.Password, app.Config.Db.Name)

	// Create connection
	db, err := sql.Open("postgres", postgresConfig)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Conexión exitosa a la BD en: " + app.Config.Db.Ip + ":" + app.Config.Db.Port + " using database " + app.Config.Db.Name)

	return db
}
