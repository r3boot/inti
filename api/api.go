package api

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/r3boot/inti/queue"
)

var listenAddr string = "localhost:7231"

var FrameQueue chan queue.FrameQueueItem

func setupRouting () (err error) {
    r := mux.NewRouter()
    r.HandleFunc("/ping", PingHandler).Methods("GET")
    r.HandleFunc("/config", ConfigHandler).Methods("GET")
    r.HandleFunc("/frame", FrameHandler).Methods("PUT")
    r.HandleFunc("/render", RenderHandler).Methods("PUT")
    r.HandleFunc("/js/{name}", FileServerHandler).Methods("GET")
    r.HandleFunc("/img/{name}", FileServerHandler).Methods("GET")
    r.HandleFunc("/css/{name}", FileServerHandler).Methods("GET")
    r.HandleFunc("/", FileServerHandler).Methods("GET")
    http.Handle("/", r)
    return
}

func Setup (listen_addr string) (err error) {
    if err = setupRouting(); err != nil { log.Fatal(err) }
    return
}

func Run () (err error) {
    log.Print("Starting server on http://"+listenAddr)
    if err = http.ListenAndServe(listenAddr, nil); err != nil {
        log.Fatal(err)
    }
    return
}
