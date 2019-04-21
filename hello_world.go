package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gorilla/mux"
	"github.com/prytoegrian/swapi/handlers"
)

func main() {
	log.Println("hello world")
	// flag log each route, operation
	var flagvar int
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	flag.Parse()
	// log.Println(flagvar)
	db := newDb()

	r := mux.NewRouter()
	h := handlers.NewHandler(db)

	r.HandleFunc("/peoples", h.AllPeoples)
	r.HandleFunc("/peoples/{id:[0-9]+}", h.OnePeople)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func newDb() *sqlite3.Conn {
	d, err := sqlite3.Open(os.Getenv("GOPATH") + "/src/github.com/prytoegrian/swapi/database/swapi.dat")
	if err != nil {
		log.Fatal(err)
	}
	d.BusyTimeout(5 * time.Second)

	return d
}

// flag
