package app

import (
	"database/sql"
	"embed"
	"web_final/config"
)

// App contains and supplies available configurations and connections
type App struct {
	Config         config.Configuration // Configuration file
	Db             *sql.DB              // Database connection
	Res            *embed.FS            // Resources from the embedded filesystem
	ScheduledTasks Scheduled            // Scheduled contains a struct of all scheduled functions
}
