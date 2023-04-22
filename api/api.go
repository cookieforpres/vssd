package api

import (
	"fmt"
	"log"
	"net/http"
	"vssd/ssd"

	"github.com/gorilla/mux"
)

type API struct {
	Host    string
	Port    int
	Verbose bool
	SSD     *ssd.SSD
}

func New(host string, port int, name string, size int, verbose bool) *API {
	return &API{
		Host:    host,
		Port:    port,
		Verbose: verbose,
		SSD:     ssd.New(name, size),
	}
}

func (a *API) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/write", a.Write)
	r.HandleFunc("/read/{name}", a.Read)
	r.HandleFunc("/delete/{name}", a.Delete)

	http.Handle("/", r)
	if a.Verbose {
		log.Printf("starting http server on %s:%d\n", a.Host, a.Port)
	}
	return http.ListenAndServe(fmt.Sprintf("%s:%d", a.Host, a.Port), nil)
}
