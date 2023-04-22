package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vssd/ssd"

	"github.com/gorilla/mux"
)

type WriteMessage struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func exist(m map[string]string, s string) bool {
	for k := range m {
		if k == s {
			return true
		}
	}

	return false
}

func (a *API) Write(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var jsonData WriteMessage
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "invalid json", "error": true}`))
		return
	}

	err = a.SSD.Write(ssd.NewNode(jsonData.Name), []byte(jsonData.Data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error", "error": true}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ok", "error": false}`))
}

func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	if ok := exist(vars, "name"); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "please provide a name", "error": true}`))
		return
	}

	data, err := a.SSD.Read(vars["name"])
	if err != nil {
		if err == ssd.ErrNoNode {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"message": "%s", "error": true}`, err.Error())))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "internal server error", "error": true}`))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	if ok := exist(vars, "name"); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.SSD.Delete(vars["name"])
	if err != nil {
		if err == ssd.ErrNoNode {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"message": "%s", "error": true}`, err.Error())))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "internal server error", "error": true}`))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ok", "error": false}`))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "pong", "error": false}`))
}
