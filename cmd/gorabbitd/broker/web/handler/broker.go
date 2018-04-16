package handler

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/broker/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/web"
	"net/http"
)

type BrokerHandler struct {
	repository repository.BrokerRepository
}

func NewBrokerHandler(repository repository.BrokerRepository) *BrokerHandler {
	return &BrokerHandler{repository}
}
func (b *BrokerHandler) Create(w http.ResponseWriter, r *http.Request) {

	web.Respond(w, nil, http.StatusOK)
}
