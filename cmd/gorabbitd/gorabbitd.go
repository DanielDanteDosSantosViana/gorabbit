package main

import (
	broker_routes "github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/broker/web/route"
	queue_routes "github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/queue/web/route"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/web"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/collector"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {

	enviroment.Load()
	session, err := db.NewSession()
	if err != nil {
		log.Panicf(err.Error())
	}

	negroniAPI := negroni.New()

	_,err = collector.NewConn()
	if err!=nil{
		log.Error(err)
	}


	api := web.NewAPI()
	broker_routes.AddAPI(session, api)
	queue_routes.AddAPI(session, api)

	negroniAPI.UseHandler(api)

	log.Info(" gorabbit running on port ", enviroment.Conf.Service.Port)

	err = http.ListenAndServe(":"+enviroment.Conf.Service.Port, negroniAPI)
	if err != nil {
		log.Fatal("Error init Server : ", err)
	}

}
