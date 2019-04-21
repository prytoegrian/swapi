package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gorilla/mux"
	"github.com/prytoegrian/swapi/repository"
)

func NewHandler(d *sqlite3.Conn) Handler {
	return Handler{
		db: d,
	}
}

type Handler struct {
	db *sqlite3.Conn
}

// AllPeoples fetches all peoples.
func (h Handler) AllPeoples(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		m, _ := json.MarshalIndent(repository.AllPeoples(h.db), "", " ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(m)
	case "POST":
		d := json.NewDecoder(r.Body)
		j := h.postPeople(d)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	case "OPTIONS":
		fallthrough
	default:
		supported := "GET, POST, OPTIONS"
		w.Header().Set("Allow", supported)
		// output json
		io.WriteString(w, "Supported methods :"+supported)
	}
}

// @TODO: apply jsend
func (h Handler) postPeople(d *json.Decoder) []byte {
	badRequest := map[string]string{"code": "400", "message": "Bad request"}
	var j map[string]string
	var p repository.People
	err := d.Decode(&p)
	if err != nil {
		log.Print(err)

		j = badRequest
	} else {
		if err := repository.PostPeople(h.db, p); err != nil {
			log.Print(err)
			j = badRequest
		} else {
			j = map[string]string{"code": "200", "message": "OK"}
		}
	}

	m, _ := json.Marshal(j)

	return m
}

// OnePeople fetches one people.
func (h Handler) OnePeople(w http.ResponseWriter, r *http.Request) {
	qs := mux.Vars(r)
	id, _ := strconv.Atoi(qs["id"])

	switch r.Method {
	case "GET":
		j := h.getPeople(id)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	case "PUT":
		d := json.NewDecoder(r.Body)
		j := h.putPeople(id, d)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	case "DELETE":
		j := h.deletePeople(id)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	case "OPTIONS":
		fallthrough
	default:
		supported := "GET, PUT, DELETE, OPTIONS"
		w.Header().Set("Allow", supported)
		// output json
		io.WriteString(w, "Supported methods :"+supported)
	}
}

// @TODO: apply jsend everywhere
func (h Handler) getPeople(id int) []byte {
	var m []byte
	people, err := repository.PeopleByID(h.db, id)
	if err != nil {
		notFound := map[string]string{"code": "404", "message": "People #" + strconv.Itoa(id) + " not found"}
		m, _ = json.Marshal(notFound)
	} else {
		m, _ = json.MarshalIndent(*people, "", " ")
	}

	return m
}

// @TODO: apply jsend
func (h Handler) putPeople(id int, d *json.Decoder) []byte {
	badRequest := map[string]string{"code": "400", "message": "Bad request"}
	var j map[string]string
	var p repository.People
	err := d.Decode(&p)
	if err != nil {
		j = badRequest
	} else {
		if err := repository.PutPeople(h.db, id, p); err != nil {
			j = badRequest
		} else {
			j = map[string]string{"code": "200", "message": "OK"}
		}
	}

	m, _ := json.Marshal(j)

	return m
}

func (h Handler) deletePeople(id int) []byte {
	notFound := map[string]string{"code": "404", "message": "People #" + strconv.Itoa(id) + " not found"}
	var j map[string]string

	if err := repository.DeletePeople(h.db, id); err != nil {
		j = notFound
	} else {
		j = map[string]string{"code": "200", "message": "OK"}
	}

	m, _ := json.Marshal(j)

	return m
}
