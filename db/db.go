package db

import (
	"database/sql"

	"github.com/lib/pq"
)

var db *sql.DB
var listener *pq.Listener

// Connect to dabasase
func Connect(url string) error {
	c, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	db = c
	return nil
}
