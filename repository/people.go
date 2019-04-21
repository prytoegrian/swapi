package repository

import (
	"errors"
	"log"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// People represents a well-formed people
type People struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Height    int        `json:"height"`
	Mass      int        `json:"mass"`
	Hair      string     `json:"hair"`
	Skin      string     `json:"skin"`
	Eye       string     `json:"eye"`
	BirthYear string     `json:"birth_year"`
	Gender    string     `json:"gender"`
	Homeworld int        `json:"homeworld"`
	Films     string     `json:"films"`
	Species   string     `json:"species"`
	Vehicles  []Vehicle  `json:"vehicles"`
	Starships []Starship `json:"starships"` //?
	Created   string     `json:"_created"`
	Edited    string     `json:"_edited"`
	URL       string     `json:"url"`
}

// AllPeoples fetches all peoples from storage
func AllPeoples(db *sqlite3.Conn) []People {
	peoples := make([]People, 0)

	stmt, err := db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
        FROM people
        ORDER BY created`)
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

		p := buildPeople(stmt)
		p.Vehicles = allVehiclesByPeopleID(db, p.ID)
		p.Starships = allStarshipsByPeopleID(db, p.ID)
		peoples = append(peoples, p)
	}

	return peoples
}

// PostPeople set one people into storage
func PostPeople(db *sqlite3.Conn, p People) error {
	l, err := lastPeople(db)
	if err != nil {
		return err
	}
	futureID := (*l).ID + 1

	now := time.Now()
	date := now.Format(time.RFC3339)
	stmt, err := db.Prepare(`INSERT INTO people
        (id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	err = stmt.Exec(
		futureID,
		p.Name,
		p.Height,
		p.Mass,
		p.Hair,
		p.Skin,
		p.Eye,
		p.BirthYear,
		p.Gender,
		p.Homeworld,
		date,
		date,
		p.URL,
	)

	if err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}

// lastPeople fetches last people from storage
func lastPeople(db *sqlite3.Conn) (*People, error) {
	stmt, err := db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
        FROM people
        ORDER BY created DESC
        LIMIT 1`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	hasRow, err := stmt.Step()
	if err != nil {
		log.Fatal("Step gave error :" + err.Error())
	}
	if !hasRow {
		return nil, errors.New("No people in storage")
	}

	p := buildPeople(stmt)
	p.Vehicles = allVehiclesByPeopleID(db, p.ID)
	p.Starships = allStarshipsByPeopleID(db, p.ID)
	return &p, nil
}

// PeopleByID fetches one people from storage
func PeopleByID(db *sqlite3.Conn, id int) (*People, error) {
	stmt, err := db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
        FROM people
        WHERE id = ?
        ORDER BY created`, id)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()

	hasRow, err := stmt.Step()
	if err != nil {
		log.Fatal("Step gave error :" + err.Error())
	}
	if !hasRow {
		return nil, errors.New("Unknown id")
	}

	p := buildPeople(stmt)
	p.Vehicles = allVehiclesByPeopleID(db, p.ID)
	p.Starships = allStarshipsByPeopleID(db, p.ID)
	return &p, nil
}

// PutPeople updates a people into storage
func PutPeople(db *sqlite3.Conn, id int, p People) error {
	now := time.Now()
	_, err := PeopleByID(db, id)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`UPDATE people
        SET name = ?, height = ?, mass = ?, hair_color = ?, skin_color = ?, eye_color = ?, birth_year = ?, gender = ?, homeworld = ?, edited = ?, url = ?
        WHERE id = ?`)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	err = stmt.Exec(
		p.Name,
		p.Height,
		p.Mass,
		p.Hair,
		p.Skin,
		p.Eye,
		p.BirthYear,
		p.Gender,
		p.Homeworld,
		now.Format(time.RFC3339),
		p.URL,
		id,
	)

	if err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}

func buildPeople(s *sqlite3.Stmt) People {
	var id int
	var name string
	var height int
	var mass int
	var hair string
	var skin string
	var eye string
	var birth_year string
	var gender string
	var homeworld int
	var created string
	var edited string
	var url string

	err := s.Scan(&id, &name, &height, &mass, &hair, &skin, &eye, &birth_year, &gender, &homeworld, &created, &edited, &url)
	if err != nil {
		log.Fatal("Scan gave error :" + err.Error())
	}

	return People{
		ID:        id,
		Name:      name,
		Height:    height,
		Mass:      mass,
		Hair:      hair,
		Skin:      skin,
		Eye:       eye,
		BirthYear: birth_year,
		Gender:    gender,
		Homeworld: homeworld,
		Created:   created,
		Edited:    edited,
		URL:       url,
	}
}

// DeletePeople unsets a people from storage
func DeletePeople(db *sqlite3.Conn, id int) error {
	_, err := PeopleByID(db, id)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`DELETE FROM people WHERE id = ?`)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	if err = stmt.Exec(id); err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}
