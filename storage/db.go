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

// Object represents an upload on the server, as it is in the database
type Object struct {
	ID           string `json:"id"`
	Date         int64  `json:"date"`
	Uploader     string `json:"uploader"`
	Type         Type   `json:"type"`
	Location     string `json:"location"`
	OriginalName string `json:"originalName"`
	MimeType     string `json:"mimeType"`
}

// Type represents the type of upload
type Type int

// File = 0, URL = 1
const (
	File Type = iota
	URL
)

// db is globalized for now
var db *sql.DB

// IDExists attempts to query one row for the passed ID.
// If it the query scans without error, that means the ID was found and this will return true.
// If the scan method returns ErrNoRows, that means it was not found and will return false.
func IDExists(ID string) bool {
	log.Printf("Attempting to query for ID %s\n...", ID)
	r := db.QueryRow("SELECT id FROM uploads WHERE id=$1", ID)
	e := r.Scan()
	if e != nil && e == sql.ErrNoRows {
		log.Printf("ID %s not found.\n", ID)
		return false
	} else if e == nil {
		log.Panicf("IDExists ID %s had an error: %v\n", ID, e)
		return false
	}
	return true
}

// AddUpload adds an uploaded file to the database.
// Returns true if database insert succeeds - false otherwise.
func AddUpload(object Object) (bool, error) {
	log.Printf("Attempting to add %v to database...", object)
	if object.ID == "" {
		log.Println("AddUpload ID field is empty!")
		return false, errors.New("ID field cannot be empty")
	}
	_, e := db.Exec(
		"INSERT INTO uploads VALUES($1, $2, $3, $4, $5, $6, $7)",
		object.ID,
		object.Date,
		object.Uploader,
		object.Type,
		object.Location,
		object.OriginalName,
		object.MimeType,
	)
	if e != nil {
		log.Panicln(e)
	}
	return true, nil
}

// GetUpload will return an object containing the basic upload data to view a file/URL.
// Object will be populated with the Type and Location of the upload, if it exists.
// An error will be returned if the lookup failed.
func GetUpload(ID string) (Object, error) {
	o := Object{}
	log.Printf("Attempting to get upload ID %s\n", ID)
	r := db.QueryRow("SELECT type, location, mimetype FROM uploads WHERE id=?", ID)
	e := r.Scan(&o.Type, &o.Location, &o.MimeType)
	if e != nil {
		log.Printf("GetUpload ID %s had an error: %v\n", ID, e)
		return Object{}, e
	}
	log.Printf("Query succeeded for ID %s (%v)\n", ID, r)
	return o, nil
}

// Initialize will create the proper table in the database, if not present.
func Initialize(d *sql.DB) {
	db = d
	_, e := d.Exec(`
		CREATE TABLE IF NOT EXISTS uploads (
			id CHAR(6),
			date BIGINT,
			uploader VARCHAR(255),
			type INT,
			location TEXT,
			original_name TEXT,
			mimetype TEXT
		);
	`)
	if e != nil {
		log.Panicln(e)
	}
}
