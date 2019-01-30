package handlers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/juju/errors"
)

// Favicon will return the favicon file, if it exists
func Favicon(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open("./assets/favicon.ico")
	if e != nil {
		log.Println(errors.Annotate(e, "Favicon"))
		w.WriteHeader(404)
		return
	}
	_, e = io.Copy(io.Writer(w), io.Reader(f))
	w.WriteHeader(200)
	return
}
