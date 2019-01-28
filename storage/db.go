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
	IDExists(ID string) (bool, error)
	AddUpload(object Object) (bool, error)
	GetUpload(ID string) (Object, error)
}

// Handler consists of the database connection.
// db can be lowercase (not exported) since it is only used here (?)
type Handler struct {
	db *sql.DB
}

// IDExists attempts to query one row for the passed ID.
// If it the query scans without error, that means the ID was found and this will return true.
// If the scan method returns ErrNoRows, that means it was not found and will return false.
func (h Handler) IDExists(ID string) (bool, error) {
	log.Printf("Attempting to query for ID %s\n...", ID)
	r := h.db.QueryRow("SELECT id FROM uploads WHERE id=?", ID)
	e := r.Scan()
	if e != nil && e == sql.ErrNoRows {
		log.Printf("ID %s not found.\n", ID)
		return false, nil
	} else if e == nil {
		log.Printf("IDExists ID %s had an error: %v\n", ID, e)
		return false, e
	}
	return false, nil
}

// AddUpload adds an uploaded file to the database.
// Returns true if database insert succeeds - false otherwise.
func (h Handler) AddUpload(object Object) (bool, error) {
	log.Printf("Attempting to add %v to database...", object)
	if object.ID == "" {
		log.Println("AddUpload ID field is empty!")
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

// GetUpload will return an object containing the basic upload data to view a file/URL.
// Object will be populated with the Type and Location of the upload, if it exists.
// An error will be returned if the lookup failed.
func (h Handler) GetUpload(ID string) (Object, error) {
	var t Type   // Type
	var l string // Location
	log.Printf("Attempting to get basic upload ID %s\n", ID)
	r := h.db.QueryRow("SELECT type, location FROM uploads WHERE id=?", ID)
	e := r.Scan(&t, &l)
	if e != nil {
		log.Printf("GetUpload ID %s had an error: %v\n", ID, e)
		return Object{}, e
	}
	log.Printf("Query succeeded for ID %s: %v\n", ID, r)
	return Object{Type: t, Location: l}, nil
}
