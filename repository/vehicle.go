package repository

import (
	"log"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// Vehicle represents a well-formed vehicle
type Vehicle struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Model                string `json:"model"`
	Manufacturer         string `json:"manufacturer"`
	CostInCredits        string `json:"cost_in_credits"`
	Length               string `json:"length"`
	MaxAtmospheringSpeed string `json:"max_atmosphering_speed"`
	Crew                 string `json:"crew"`
	Passengers           string `json:"passengers"`
	CargoCapacity        string `json:"cargo_capacity"`
	Consumables          string `json:"consumables"`
	VehicleClass         string `json:"vehicle_class"`
	Pilots               string `json:"pilots"`
	Films                string `json:"films"`
	Created              string `json:"_created"`
	Edited               string `json:"_edited"`
	URL                  string `json:"url"`
}

func allVehiclesByPeopleID(db *sqlite3.Conn, id int) []Vehicle {
	vs := make([]Vehicle, 0)

	stmt, err := db.Prepare(`SELECT id, name, model, manufacturer, cost_in_credits, length, max_atmosphering_speed, crew, passengers, cargo_capacity, consumables, vehicle_class, created, edited, url
        FROM people_vehicles pv
            INNER JOIN vehicles v ON pv.vehicles = v.id
        WHERE pv.people = ?`, id)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			log.Fatal("Step gave error :" + err.Error())
		}
		if !hasRow {
			break
		}

		// improvement : mass fetching of vehicle for all id people
		v := buildVehicle(stmt)
		vs = append(vs, v)
	}

	return vs
}

func buildVehicle(s *sqlite3.Stmt) Vehicle {
	var id int
	var name string
	var model string
	var manufacturer string
	var cost_in_credits string
	var length string
	var max_atmosphering_speed string
	var crew string
	var passengers string
	var cargo_capacity string
	var consumables string
	var vehicle_class string
	var created string
	var edited string
	var url string

	err := s.Scan(&id, &name, &model, &manufacturer, &cost_in_credits, &length, &max_atmosphering_speed, &crew, &passengers, &cargo_capacity, &consumables, &vehicle_class, &created, &edited, &url)
	if err != nil {
		log.Fatal("Scan gave error :" + err.Error())
	}

	return Vehicle{
		ID:                   id,
		Name:                 name,
		Model:                model,
		Manufacturer:         manufacturer,
		CostInCredits:        cost_in_credits,
		Length:               length,
		MaxAtmospheringSpeed: max_atmosphering_speed,
		Crew:                 crew,
		Passengers:           passengers,
		CargoCapacity:        cargo_capacity,
		Consumables:          consumables,
		VehicleClass:         vehicle_class,
		Created:              created,
		Edited:               edited,
		URL:                  url,
	}
}
