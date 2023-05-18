package app

import (
	"database/sql"
	"embed"
	"web_final/config"
)

type App struct {
	Config         config.Configuration
	Db             *sql.DB
	Res            *embed.FS
	ScheduledTasks Scheduled
}
