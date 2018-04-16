package main

import (
	broker_routes "github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/broker/web/route"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"os"
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

	brokerAPI := broker_routes.API(session)

	negroniAPI.UseHandler(brokerAPI)

	log.Info(" gorabbit running on port %s ", enviroment.Conf.Service.Port)

	err = http.ListenAndServe(":"+enviroment.Conf.Service.Port, negroniAPI)
	if err != nil {
		log.Fatal("Error init Server : ", err)
	}

}
