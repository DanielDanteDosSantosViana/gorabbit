package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

type InvalidError []Invalid

func (err InvalidError) Error() string {
	var str string
	for _, v := range err {
		str = fmt.Sprintf("%s,{%s:%s}", str, v.Fld, v.Err)
	}
	return str
}

type JSONError struct {
	Error  string       `json:"error"`
	Fields InvalidError `json:"fields,omitempty"`
}

func RespondError(w http.ResponseWriter, err error, code int) {
	Respond(w, JSONError{Error: err.Error()}, code)
}

func Respond(w http.ResponseWriter, data interface{}, code int) {

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("%s : Respond %v Marshalling JSON response\n", err, jsonData)
		jsonData = []byte("{}")
	}
	io.WriteString(w, string(jsonData))
}
