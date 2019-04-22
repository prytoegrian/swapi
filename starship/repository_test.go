package starship

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

func TestAllStarshipsByPeopleIDOK(t *testing.T) {
	step = 0
	ss := repo.AllStarshipsByPeopleID(88)
	if len(ss) != 2 {
		t.Error("No starship for this people")
	}
}

func TestAllStarshipsByPeopleIDKO(t *testing.T) {
	step = 3
	ss := repo.AllStarshipsByPeopleID(88)
	if len(ss) != 0 {
		t.Error("There's starship for this people")
	}
}
