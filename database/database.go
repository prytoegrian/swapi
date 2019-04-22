package database

import (
	"log"
	"os"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// Database describes accesses to a storage
type Database interface {
	Prepare(string, ...interface{}) (Stmt, error)
}

// Db represents a connection to the storage
type Db struct {
	sqlite *sqlite3.Conn
}

// Prepare encapsulates the inner connection for testability
func (d Db) Prepare(sql string, args ...interface{}) (Stmt, error) {
	return d.sqlite.Prepare(sql, args...)
}

// Stmt represents a query statement
// Cf. sqlite3.Stmt
type Stmt interface {
	Close() error
	Step() (bool, error)
	Exec(...interface{}) error
	Scan(dst ...interface{}) error
}

// NewDb initialise a new connection
func NewDb() Database {
	s, err := sqlite3.Open(os.Getenv("GOPATH") + "/src/github.com/prytoegrian/swapi/database/swapi.dat")
	if err != nil {
		log.Fatal(err)
	}
	s.BusyTimeout(5 * time.Second)

	return Db{
		sqlite: s,
	}
}
