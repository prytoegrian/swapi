package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prytoegrian/swapi/database"
	"github.com/prytoegrian/swapi/handlers"
	"github.com/prytoegrian/swapi/people"
)

func main() {
	log.Println("Le serveur écoute désormais à http://localhost:8080")
	log.Println("Pour couper le serveur, tapez simplement Ctrl-C")
	// flag log each route, operation
	var debug int
	flag.IntVar(&debug, "debug", 0, "Enable ou disable full log")
	flag.Parse()

	db := database.NewDb()
	r := mux.NewRouter()
	repo := people.NewRepo(db)
	h := handlers.NewHandler(repo)

	r.HandleFunc("/peoples", h.AllPeoples)
	r.HandleFunc("/peoples/{id:[0-9]+}", h.OnePeople)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
