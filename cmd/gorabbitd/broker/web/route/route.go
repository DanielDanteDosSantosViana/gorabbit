package route

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/broker/web/handler"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/gorilla/mux"
	"net/http"
)


func API(sessiondb db.Session) *mux.Router {

	repository := repository.NewMongoBrokerRepository(sessiondb)
	brokerHandler := handler.NewBrokerHandler(repository)

	r := mux.NewRouter().StrictSlash(true)

	api := r.PathPrefix("/v1").Subrouter()

	api.HandleFunc("/brokers", brokerHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/brokers/{id}", brokerHandler.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/brokers", brokerHandler.List).Methods(http.MethodGet)

	return api
}
