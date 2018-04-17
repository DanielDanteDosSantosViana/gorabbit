package handler

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/queue/repository"
	broker_repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/web"
	"net/http"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/queue"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type QueueHandler struct {
	repository repository.QueueRepository
	brokerRepository broker_repo.BrokerRepository
}

func NewQueueHandler(repository repository.QueueRepository, brokerRepo broker_repo.BrokerRepository) *QueueHandler {
	return &QueueHandler{repository, brokerRepo}
}

func (q *QueueHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)

	broker_id_req := vars["broker_id"]
	broker_id := bson.ObjectIdHex(broker_id_req)

	if broker, err:= q.brokerRepository.Get(broker_id);err!=nil{
		err := web.IsRequestValid(broker)
		if err !=nil{
			web.RespondError(w, err, http.StatusBadRequest)
			return
		}
	}

	queue := &queue.Queue{}

	if err := json.Unmarshal(body, queue); err != nil {
		log.WithFields(log.Fields{"queue": queue, "err": err.Error()}).Error("Payload invalid of queue request")
		web.RespondError(w, err, http.StatusUnprocessableEntity)
		return
	}

	queue.BrokerID = broker_id
	log.Println(queue.BrokerID)
	log.Println(broker_id)
	err := web.IsRequestValid(queue)
	if err !=nil{
		web.RespondError(w, err, http.StatusBadRequest)
		return
	}

	q.repository.Store(queue)

	web.Respond(w, queue, http.StatusCreated)
}

func (q *QueueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	broker_id_req := vars["broker_id"]
	idReq := vars["id"]

	broker_id := bson.ObjectIdHex(broker_id_req)
	if broker, err:= q.brokerRepository.Get(broker_id);err!=nil{
		err := web.IsRequestValid(broker)
		if err !=nil{
			web.RespondError(w, err, http.StatusBadRequest)
			return
		}
	}

	if idReq==""{
		web.RespondError(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	id := bson.ObjectIdHex(idReq)
	if err := q.repository.Delete(id);err !=nil{
		web.Respond(w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}

func (q *QueueHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	broker_id := vars["broker_id"]
	id := bson.ObjectIdHex(broker_id)

	if broker, err:= q.brokerRepository.Get(id);err!=nil{
		err := web.IsRequestValid(broker)
		if err !=nil{
			web.RespondError(w, err, http.StatusBadRequest)
			return
		}
	}

	queues, err := q.repository.ListByBrokerID(id)
	if err !=nil{
		web.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(w, queues, http.StatusOK)
}

