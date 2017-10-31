package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/symflower/sessions/2017/socrates-linz-microservice-calculator"
)

func main() {
	r := mux.NewRouter()

	var lock sync.Mutex

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()

		var op socra.Add
		if err := socra.Decode(r, &op); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		socra.Encode(w, socra.Result{op.A + op.B})
	}).Methods("POST")

	socra.RunServer(r)
}
