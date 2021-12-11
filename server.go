package main

import (
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "net/http"
)

func StartServer() {
    r := mux.NewRouter()
    RegisterEndpoints(r)
    log.Info("Starting http server")

    // TODO: make port configurable
    const port = ":5555"
    log.Fatal(http.ListenAndServe(port, r))
}

func ServerHandler(w http.ResponseWriter, r *http.Request) {
    log.Info("Got request from: ", r.RemoteAddr)
}

func RegisterEndpoints(r *mux.Router) {
    r.HandleFunc("/", ServerHandler).Methods("GET").Schemes("http")
    http.Handle("/", r)
}
