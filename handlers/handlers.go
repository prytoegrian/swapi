package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prytoegrian/swapi/people"
)

// NewHandler initialise a new handler
func NewHandler(r people.Repository) Handler {
	return Handler{
		r: r,
	}
}

// Handler contains all routes descriptions
type Handler struct {
	r people.Repository
}

// AllPeoples work on all peoples.
func (h Handler) AllPeoples(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m []byte

	switch r.Method {
	case "GET":
		filled := filledOK(h.r.AllPeoples())
		m, _ = json.MarshalIndent(filled, "", " ")
	case "POST":
		d := json.NewDecoder(r.Body)
		m = h.postPeople(d)
	case "OPTIONS":
		fallthrough
	default:
		supported := "GET, POST, OPTIONS"
		w.Header().Set("Allow", supported)
		m, _ = json.MarshalIndent(notAllowed(supported), "", " ")
	}

	w.Write(m)
}

func (h Handler) postPeople(d *json.Decoder) []byte {
	badRequest := badRequest()
	var o Output
	var p people.People
	err := d.Decode(&p)
	if err != nil {
		log.Print(err)

		o = badRequest
	} else {
		if id := h.r.PostPeople(p); id == 0 {
			o = badRequest
		} else {
			o = voidOK()
		}
	}
	m, _ := json.MarshalIndent(o, "", " ")

	return m
}

// OnePeople work on one people.
func (h Handler) OnePeople(w http.ResponseWriter, r *http.Request) {
	qs := mux.Vars(r)
	id, _ := strconv.Atoi(qs["id"])
	w.Header().Set("Content-Type", "application/json")
	var m []byte

	switch r.Method {
	case "GET":
		m = h.getPeople(id)
	case "PUT":
		d := json.NewDecoder(r.Body)
		m = h.putPeople(id, d)
	case "DELETE":
		m = h.deletePeople(id)
	case "OPTIONS":
		fallthrough
	default:
		supported := "GET, PUT, DELETE, OPTIONS"
		w.Header().Set("Allow", supported)
		m, _ = json.MarshalIndent(notAllowed(supported), "", " ")
	}
	w.Write(m)
}

func (h Handler) getPeople(id int) []byte {
	var o interface{}
	people, err := h.r.PeopleByID(id)
	if err != nil {
		o = notFound(id)
	} else {
		o = filledOK(people)
	}
	m, _ := json.MarshalIndent(o, "", " ")

	return m
}

func (h Handler) putPeople(id int, d *json.Decoder) []byte {
	badRequest := badRequest()
	var o interface{}
	var p people.People
	err := d.Decode(&p)
	if err != nil {
		o = badRequest
	} else {
		if err := h.r.PutPeople(id, p); err != nil {
			o = badRequest
		} else {
			o = voidOK()
		}
	}

	m, _ := json.MarshalIndent(o, "", " ")

	return m
}

func (h Handler) deletePeople(id int) []byte {
	var j interface{}
	if err := h.r.DeletePeople(id); err != nil {
		j = notFound(id)
	} else {
		j = voidOK()
	}

	m, _ := json.MarshalIndent(j, "", " ")

	return m
}

func voidOK() Output {
	return Output{
		Code:    200,
		Status:  "OK",
		Message: "",
	}
}

func filledOK(d interface{}) interface{} {
	ok := voidOK()
	filled := struct {
		Code    int         `json:"code"`
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    ok.Code,
		Status:  ok.Status,
		Message: ok.Message,
		Data:    d,
	}

	return filled
}

func notFound(id int) Output {
	return Output{
		Code:    404,
		Status:  "Fail",
		Message: "People #" + strconv.Itoa(id) + " not found",
	}
}

func badRequest() Output {
	return Output{
		Code:    400,
		Status:  "Fail",
		Message: "Bad request",
	}
}

func notAllowed(s string) Output {
	return Output{
		Code:    405,
		Status:  "Fail",
		Message: "Supported methode : " + s,
	}
}

// Output represents an API output
type Output struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
