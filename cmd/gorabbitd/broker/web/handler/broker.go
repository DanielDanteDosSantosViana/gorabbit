package handler

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	queue_repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/queue/repository"

	"encoding/json"
	"errors"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/web"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/collector"
)

type BrokerHandler struct {
	repository repository.BrokerRepository
	queueRepo  queue_repo.QueueRepository
	collector  *collector.Collector
}

func NewBrokerHandler(repository repository.BrokerRepository, queueRepo queue_repo.QueueRepository) *BrokerHandler {
	return &BrokerHandler{repository, queueRepo}
}

func (b *BrokerHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	broker := &broker.Broker{}

	if err := json.Unmarshal(body, broker); err != nil {
		log.WithFields(log.Fields{"broker": broker, "err": err.Error()}).Error("Payload invalid of broker request")
		web.RespondError(w, err, http.StatusUnprocessableEntity)
		return
	}

	err := web.IsRequestValid(broker)
	if err != nil {
		web.RespondError(w, err, http.StatusBadRequest)
		return
	}

	b.repository.Store(broker)
	log.WithFields(log.Fields{"broker": broker}).Info("saved with sucess")

	web.Respond(w, broker, http.StatusCreated)
}

func (b *BrokerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idReq := vars["id"]
	if idReq == "" {
		err := errors.New("id is required")
		log.WithFields(log.Fields{"id": idReq}).Error(err.Error())
		web.RespondError(w, err, http.StatusBadRequest)
		return
	}
	id := bson.ObjectIdHex(idReq)
	if err := b.repository.Delete(id); err != nil {
		log.WithFields(log.Fields{"id": idReq}).Error(err.Error())
		web.Respond(w, err, http.StatusInternalServerError)
		return
	}

	b.queueRepo.DeleteByBrokerID(id)

	log.WithFields(log.Fields{"broker with id ": idReq}).Info("deleted with sucess")
	web.Respond(w, nil, http.StatusNoContent)
}

func (b *BrokerHandler) List(w http.ResponseWriter, r *http.Request) {
	brokers, err := b.repository.List()
	if err != nil {
		log.WithFields(log.Fields{"brokers": brokers}).Error(err.Error())
		web.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(w, brokers, http.StatusOK)
}
