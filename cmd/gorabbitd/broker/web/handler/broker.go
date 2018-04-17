package handler

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/web"
	"net/http"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type BrokerHandler struct {
	repository repository.BrokerRepository
}

func NewBrokerHandler(repository repository.BrokerRepository) *BrokerHandler {
	return &BrokerHandler{repository}
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
	if err !=nil{
		web.RespondError(w, err, http.StatusBadRequest)
		return
	}

	b.repository.Store(broker)

	web.Respond(w, broker, http.StatusCreated)
}

func (b *BrokerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idReq := vars["id"]
	if idReq==""{
		web.RespondError(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	id := bson.ObjectIdHex(idReq)
	if err := b.repository.Delete(id);err !=nil{
		web.Respond(w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}

func (b *BrokerHandler) List(w http.ResponseWriter, r *http.Request) {
	brokers, err := b.repository.List()
	if err !=nil{
		web.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(w, brokers, http.StatusOK)
}

