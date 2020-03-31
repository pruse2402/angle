package handlers

import (
	"angle/src/errs"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func parseJson(w http.ResponseWriter, params io.ReadCloser, data interface{}) bool {
	if params != nil {
		defer params.Close()
	}

	err := json.NewDecoder(params).Decode(data)
	if err == nil {
		return true
	}

	e := &errs.AppError{
		Message: "Invalid JSON.",
		Err:     err.Error(),
	}

	respondError(w, http.StatusBadRequest, e)
	return false
}

func respondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func renderError(w http.ResponseWriter, err interface{}) {
	var status int
	switch e := err.(type) {
	case *errs.UIErr:
		status = e.Code
	case *errs.ApiErr:
		status = e.StatusCode
	case []*errs.ApiErr:
		status = e[0].StatusCode
		err = struct {
			Message string         `json:"message"`
			Errors  []*errs.ApiErr `json:"errors"`
		}{
			e[0].Message,
			e,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Printf("ERROR: renderError - %q\n", err)
	}
	return
}

func renderJson(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// We don't have to write body, If status code is 204 (No Content)
	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func isObjectIDValid(id string) bool {
	return bson.IsObjectIdHex(id)
}
