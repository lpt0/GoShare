// Package storage provides related functions.
// This contains the struct for handling the database object, as well as executing pre-defined queries.
package storage

import (
	"database/sql"
	"errors"
	"log"
)

// Storage interface describes the methods available to interact with the database
type Storage interface {
	IDExists(ID string) bool
	AddUpload(data map[string]string) bool
}

// Handler consists of the database connection.
// db can be lowercase (not exported) since it is only used here (?)
type Handler struct {
	db *sql.DB
}

// IDExists attempts to query one row for the passed ID.
// If it the query scans without error, that means the ID was found and this will return true.
// If the scan method returns ErrNoRows, that means it was not found and will return false.
func (h Handler) IDExists(ID string) bool {
	log.Printf("Attempting to query for ID %s\n...", ID)
	r := h.db.QueryRow("SELECT id FROM uploads WHERE id=?", ID)
	if r.Scan() == sql.ErrNoRows {
		log.Printf("ID %s not found.\n", ID)
	}
	return true
}

// AddUpload adds an uploaded file to the database.
// Returns true if database insert succeeds - false otherwise.
func (h Handler) AddUpload(object Object) (bool, error) {
	log.Printf("Attempting to add %v to database...", object)
	if object.ID == "" {
		log.Println("ID field is empty!")
		return false, errors.New("ID field cannot be empty")
	}
	_, e := h.db.Exec(
		"INSERT INTO uploads VALUES(?, ?, ?, ?, ?, ?)",
		object.ID,
		object.Date,
		object.Uploader,
		object.Type,
		object.Location,
		object.OriginalName,
	)
	if e != nil {
		log.Panicln(e)
	}
	return true, nil
}
