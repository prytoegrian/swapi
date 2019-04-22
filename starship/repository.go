package starship

import (
	"log"

	d "github.com/prytoegrian/swapi/database"
)

// NewRepo initialises a new starship repository
func NewRepo(db d.Database) Repository {
	return Repository{
		db: db,
	}
}

// Repository is a starship repository
type Repository struct {
	db d.Database
}

// Starship represents a well-formed starship
type Starship struct {
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
	HyperdriveRating     string `json:"hyperdrive_rating"`
	MGLT                 string `json:"mglt"`
	StarshipClass        string `json:"starship_class"`
	Pilots               string `json:"pilots"`
	Films                string `json:"films"`
	Created              string `json:"_created"`
	Edited               string `json:"_edited"`
	URL                  string `json:"url"`
}

// AllStarshipsByPeopleID get all starships associated to a people
func (r Repository) AllStarshipsByPeopleID(id int) []Starship {
	ss := make([]Starship, 0)

	stmt, err := r.db.Prepare(`SELECT id, name, model, manufacturer, cost_in_credits, length, max_atmosphering_speed, crew, passengers, cargo_capacity, consumables, hyperdrive_rating, mglt, starship_class, created, edited, url
        FROM people_starships ps
            INNER JOIN starships s ON ps.starships = s.id
        WHERE ps.people = ?`, id)
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
			// The query is finished
			break
		}
		s := buildStarship(stmt)
		ss = append(ss, s)
	}

	return ss
}

func buildStarship(stmt d.Stmt) Starship {
	// Use Scan to access column data from a row
	var id int
	var name string
	var model string
	var manufacturer string
	var costInCredits string
	var length string
	var maxAtmospheringSpeed string
	var crew string
	var passengers string
	var cargoCapacity string
	var consumables string
	var hyperdriveRating string
	var mglt string
	var starshipClass string
	var created string
	var edited string
	var url string

	err := stmt.Scan(&id, &name, &model, &manufacturer, &costInCredits, &length, &maxAtmospheringSpeed, &crew, &passengers, &cargoCapacity, &consumables, &hyperdriveRating, &mglt, &starshipClass, &created, &edited, &url)
	if err != nil {
		log.Fatal("Scan gave error :" + err.Error())
	}
	// improvement : mass fetching of starships for all id people
	return Starship{
		ID:                   id,
		Name:                 name,
		Model:                model,
		Manufacturer:         manufacturer,
		CostInCredits:        costInCredits,
		Length:               length,
		MaxAtmospheringSpeed: maxAtmospheringSpeed,
		Crew:                 crew,
		Passengers:           passengers,
		CargoCapacity:        cargoCapacity,
		Consumables:          consumables,
		HyperdriveRating:     hyperdriveRating,
		MGLT:                 mglt,
		StarshipClass:        starshipClass,
		Created:              created,
		Edited:               edited,
		URL:                  url,
	}
}
