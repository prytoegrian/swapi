package vehicle

import (
	"testing"

	"github.com/prytoegrian/swapi/database"
)

type DataDouble struct{}

type StmtDouble struct{}

func (d DataDouble) Prepare(sql string, args ...interface{}) (database.Stmt, error) {
	return StmtDouble{}, nil
}

func (s StmtDouble) Close() error {
	return nil
}

var step int

func (s StmtDouble) Step() (bool, error) {
	step++
	return (step <= 2), nil
}

func (s StmtDouble) Exec(...interface{}) error {
	return nil
}

func (s StmtDouble) Scan(dst ...interface{}) error {
	return nil
}

var repo = NewRepo(DataDouble{})

func TestAllVehiclesByPeopleIDOK(t *testing.T) {
	step = 0
	vs := repo.AllVehiclesByPeopleID(88)
	if len(vs) != 2 {
		t.Error("No vehicle for this people")
	}
}

func TestAllVehiclesByPeopleIDKO(t *testing.T) {
	step = 3
	vs := repo.AllVehiclesByPeopleID(88)
	if len(vs) != 0 {
		t.Error("There's vehicle for this people")
	}
}
