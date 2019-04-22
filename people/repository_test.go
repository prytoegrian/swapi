package people

import (
	"errors"
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

var exec error

func (s StmtDouble) Exec(...interface{}) error {
	return exec
}

func (s StmtDouble) Scan(dst ...interface{}) error {
	return nil
}

var repo = NewRepo(DataDouble{})

func TestAllPeoplesOK(t *testing.T) {
	step = 1
	ps := repo.AllPeoples()
	if len(ps) != 1 {
		t.Error("No people")
	}
}

func TestAllPeoplesKO(t *testing.T) {
	step = 3
	ps := repo.AllPeoples()
	if len(ps) != 0 {
		t.Error("There's people")
	}
}

func TestPostPeopleNoLast(t *testing.T) {
	step = 2
	p := People{
		ID:   97,
		Name: "Boba Fett",
	}
	if id := repo.PostPeople(p); id != 1 {
		t.Error("There's people in storage")
	}
}

func TestPeopleByIDOK(t *testing.T) {
	step = 1
	if _, err := repo.PeopleByID(5); err != nil {
		t.Error("There's no people with this id")
	}
}

func TestPeopleByIDKO(t *testing.T) {
	step = 2
	if _, err := repo.PeopleByID(7); err == nil {
		t.Error("There's people with this id")
	}
}

func TestPutPeopleNoPeople(t *testing.T) {
	step = 2
	p := People{
		ID:   874,
		Name: "Jango Fett",
	}
	if err := repo.PutPeople(6, p); err == nil {
		t.Error("Found people with this id")
	}
}

func TestPutPeopleFail(t *testing.T) {
	step = 1
	exec = errors.New("")
	p := People{
		ID:   874,
		Name: "Jango Fett",
	}
	if err := repo.PutPeople(6, p); err == nil {
		t.Error("Fail exec")
	}
}

func TestPutPeopleOK(t *testing.T) {
	step = 1
	exec = nil
	p := People{
		ID:   874,
		Name: "Jango Fett",
	}
	if err := repo.PutPeople(6, p); err != nil {
		t.Error("Put failed")
	}
}

func TestDeletePeopleNoPeople(t *testing.T) {
	step = 2
	if err := repo.DeletePeople(15); err == nil {
		t.Error("Found people with this id")
	}
}

func TestDeletePeopleFail(t *testing.T) {
	step = 1
	exec = errors.New("")

	if err := repo.DeletePeople(18); err == nil {
		t.Error("Fail exec")
	}
}

func TestDeletePeopleOK(t *testing.T) {
	step = 1
	exec = nil
	if err := repo.DeletePeople(21); err != nil {
		t.Error("Delete failed")
	}
}
