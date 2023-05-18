package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_final/app"
	"web_final/config"
	"web_final/database"
	"web_final/models"
	"web_final/routes"
)

var res embed.FS

func main() {

	appLoaded := app.App{}

	appLoaded.Config = config.LoadConfig()

	appLoaded.Res = &res

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Println("No se pudo crear el log directorio")
			log.Println(err)
			return
		}
	}

	file, err := os.OpenFile("logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)

	appLoaded.Db = database.ConnectDB(&appLoaded)
	if appLoaded.Config.Db.AutoMigrate {
		err = models.RunAllMigrations(&appLoaded)
		if err != nil {
			log.Println(err)
			return
		}
	}

	appLoaded.ScheduledTasks = app.Scheduled{
		EveryReboot: []func(app *app.App){models.ScheduledSessionCleanup},
		EveryMinute: []func(app *app.App){models.ScheduledSessionCleanup},
	}

	routes.GetRoutes(&appLoaded)
	routes.PostRoutes(&appLoaded)

	server := &http.Server{Addr: appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port}
	go func() {
		log.Println("Servidor Iniciado" + appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("No se pudo inicializar servidor %s: %v\n", appLoaded.Config.Listen.Ip+":"+appLoaded.Config.Listen.Port, err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	stop := make(chan struct{})
	go app.RunScheduledTasks(&appLoaded, 100, stop)

	<-interrupt
	log.Println("Se obtuvo una señal de interrupción, apagando servidor...")

	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("No se pudo cerrar el servidor: %v\n", err)
	}
}
