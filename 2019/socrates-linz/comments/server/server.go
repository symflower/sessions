package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // Include SQLite as SQL driver.

	"github.com/symflower/sessions/2019/socrates-linz/comments/controller"
	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

// Run starts our application server.
func Run() {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:unique%d?mode=memory&cache=shared", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}

	err = model.CommentInit(db)
	if err != nil {
		panic(err)
	}
	err = model.UserInit(db)
	if err != nil {
		panic(err)
	}

	middleware := func(handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request %s with %s", r.URL.Path, r.Method)

			handler(db, w, r)

			log.Printf("Done")
		}
	}

	http.HandleFunc("/", middleware(controller.HandleIndex))
	http.HandleFunc("/register", middleware(controller.HandleRegister))

	log.Printf("Listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
