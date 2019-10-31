package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	var list = map[int]string{}
	var lockList sync.Mutex

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		lockList.Lock()
		defer lockList.Unlock()

		id := len(list) + 1
		list[id] = time.Now().String()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(id)
		if err != nil {
			panic(err)
		}
	}).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lockList.Lock()
		defer lockList.Unlock()

		// FIXME very expensive computation!
		time.Sleep(5 * time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(list)
		if err != nil {
			panic(err)
		}
	}).Methods("GET")

	runServer(r)
}

func runServer(r *mux.Router) {
	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())

	n.UseHandler(r)

	n.Run(":8080")
}
