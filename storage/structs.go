// Package storage provides related functions
package storage

// Object represents an upload on the server, as it is in the database
type Object struct {
	ID           string `json:"id"`
	Date         int64  `json:"date"`
	Uploader     string `json:"uploader"`
	Type         Type   `json:"type"`
	Location     string `json:"location"`
	OriginalName string `json:"originalName"`
}

// Type represents the type of upload
type Type int

// File = 0, URL = 1
const (
	File Type = iota
	URL
)
