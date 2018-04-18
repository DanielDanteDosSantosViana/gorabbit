package route

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/queue/web/handler"
	broker_repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/queue/repository"

	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/gorilla/mux"
	"net/http"
)

func AddAPI(sessiondb db.Session, api *mux.Router) *mux.Router {

	repository := repository.NewMongoQueueRepository(sessiondb)
	brokerRepo := broker_repo.NewMongoBrokerRepository(sessiondb)

	queueHandler := handler.NewQueueHandler(repository, brokerRepo)

	api.HandleFunc("/brokers/{broker_id}/queues", queueHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/brokers/{broker_id}/queues/{id}", queueHandler.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/brokers/{broker_id}/queues", queueHandler.List).Methods(http.MethodGet)

	return api
}
