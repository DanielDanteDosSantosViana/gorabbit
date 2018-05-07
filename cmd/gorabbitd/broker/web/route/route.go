package route

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd/broker/web/handler"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	queue_repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/queue/repository"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/collector"
)

func AddAPI(sessiondb db.Session, api *mux.Router, collector *collector.Collector) *mux.Router {
	queueRepo := queue_repo.NewMongoQueueRepository(sessiondb)

	repository := repository.NewMongoBrokerRepository(sessiondb)
	brokerHandler := handler.NewBrokerHandler(repository, queueRepo,collector)

	api.HandleFunc("/brokers", brokerHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/brokers/{id}", brokerHandler.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/brokers", brokerHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/brokers/{id}/cmd", brokerHandler.Command).Methods(http.MethodPut)

	return api
}
