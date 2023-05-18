package models

import (
	"time"
	"web_final/app"
	"web_final/database"
)

func RunAllMigrations(app *app.App) error {

	user := User{
		Id:        1,
		Username:  "migrate",
		Password:  "migrate",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := database.Migrate(app, user)
	if err != nil {
		return err
	}

	session := Session{
		Id:         1,
		UserId:     1,
		AuthToken:  "migrate",
		RememberMe: false,
		CreatedAt:  time.Now(),
	}
	err = database.Migrate(app, session)
	if err != nil {
		return err
	}

	return nil
}
