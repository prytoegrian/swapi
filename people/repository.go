package people

import (
	"errors"
	"log"
	"time"

	d "github.com/prytoegrian/swapi/database"
	"github.com/prytoegrian/swapi/starship"
	"github.com/prytoegrian/swapi/vehicle"
)

// NewRepo initialises a new people repository
func NewRepo(db d.Database) Repository {
	return Repository{
		db: db,
	}
}

// Repository is a people repository
type Repository struct {
	db d.Database
}

// People represents a well-formed people
type People struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Height    int                 `json:"height"`
	Mass      int                 `json:"mass"`
	Hair      string              `json:"hair"`
	Skin      string              `json:"skin"`
	Eye       string              `json:"eye"`
	BirthYear string              `json:"birth_year"`
	Gender    string              `json:"gender"`
	Homeworld int                 `json:"homeworld"`
	Films     string              `json:"films"`
	Species   string              `json:"species"`
	Vehicles  []vehicle.Vehicle   `json:"vehicles"`
	Starships []starship.Starship `json:"starships"`
	Created   string              `json:"_created"`
	Edited    string              `json:"_edited"`
	URL       string              `json:"url"`
}

// AllPeoples fetches all peoples from storage
func (r Repository) AllPeoples() []People {
	peoples := make([]People, 0)

	stmt, err := r.db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
        FROM people
        ORDER BY created`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
	}
	defer stmt.Close()
	v := vehicle.NewRepo(r.db)
	s := starship.NewRepo(r.db)

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			log.Fatal("Step gave error :" + err.Error())
		}
		if !hasRow {
			break
		}

		p := buildPeople(stmt)
		p.Vehicles = v.AllVehiclesByPeopleID(p.ID)
		p.Starships = s.AllStarshipsByPeopleID(p.ID)
		peoples = append(peoples, p)
	}

	return peoples
}

// PostPeople set one people into storage
func (r Repository) PostPeople(p People) int {
	l, err := r.lastPeople()
	var futureID int
	if err != nil {
		futureID = 1
	} else {
		futureID = (*l).ID + 1
	}

	now := time.Now()
	date := now.Format(time.RFC3339)
	stmt, err := r.db.Prepare(`INSERT INTO people
        (id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal("Malformed SQL :" + err.Error())
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
		log.Fatal("Failed to exec SQL :" + err.Error())
	}

	return futureID
}

// lastPeople fetches last people from storage
func (r Repository) lastPeople() (*People, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
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
	return &p, nil
}

// PeopleByID fetches one people from storage
func (r Repository) PeopleByID(id int) (*People, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, homeworld, created, edited, url
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
	v := vehicle.NewRepo(r.db)
	s := starship.NewRepo(r.db)

	p := buildPeople(stmt)
	p.Vehicles = v.AllVehiclesByPeopleID(p.ID)
	p.Starships = s.AllStarshipsByPeopleID(p.ID)
	return &p, nil
}

// PutPeople updates a people into storage
func (r Repository) PutPeople(id int, p People) error {
	now := time.Now()
	_, err := r.PeopleByID(id)
	if err != nil {
		return err
	}

	stmt, err := r.db.Prepare(`UPDATE people
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

func buildPeople(s d.Stmt) People {
	var id int
	var name string
	var height int
	var mass int
	var hair string
	var skin string
	var eye string
	var birthYear string
	var gender string
	var homeworld int
	var created string
	var edited string
	var url string

	err := s.Scan(&id, &name, &height, &mass, &hair, &skin, &eye, &birthYear, &gender, &homeworld, &created, &edited, &url)
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
		BirthYear: birthYear,
		Gender:    gender,
		Homeworld: homeworld,
		Created:   created,
		Edited:    edited,
		URL:       url,
	}
}

// DeletePeople unsets a people from storage
func (r Repository) DeletePeople(id int) error {
	_, err := r.PeopleByID(id)
	if err != nil {
		return err
	}

	stmt, err := r.db.Prepare(`DELETE FROM people WHERE id = ?`)
	if err != nil {
		return errors.New("Failed to prepare :" + err.Error())
	}
	defer stmt.Close()

	if err = stmt.Exec(id); err != nil {
		return errors.New("Failed to exec SQL :" + err.Error())
	}

	return nil
}
