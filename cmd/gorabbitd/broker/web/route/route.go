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
	handler := handler.NewBrokerHandler(repository)

	r := mux.NewRouter().StrictSlash(true)

	api := r.PathPrefix("/v1").Subrouter()

	api.HandleFunc("/broker", handler.Create).Methods(http.MethodPost)

	return api
}
