package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func database() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS list(id SERIAL PRIMARY KEY, created TEXT)")
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	r := mux.NewRouter()

	var lockList sync.Mutex

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		lockList.Lock()
		defer lockList.Unlock()

		db := database()
		defer func() {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}()

		var lastInsertID int
		err := db.QueryRow("INSERT INTO list(created) VALUES($1) RETURNING id", time.Now().String()).Scan(&lastInsertID)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(lastInsertID)
		if err != nil {
			panic(err)
		}
	}).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lockList.Lock()
		defer lockList.Unlock()

		// FIXME very expensive computation!
		time.Sleep(5 * time.Second)

		db := database()
		defer func() {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}()

		var list = map[int]string{}
		rows, err := db.Query("SELECT id, created FROM list")
		if err != nil {
			panic(err)
		}
		defer func() {
			err := rows.Close()
			if err != nil {
				panic(err)
			}
		}()

		for rows.Next() {
			var id int
			var created string

			err := rows.Scan(&id, &created)
			if err != nil {
				panic(err)
			}

			list[id] = created
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(list)
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
