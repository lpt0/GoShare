package handlers

import (
	"goshare/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/juju/errors"
)

// Flash is the HTTP route for flash viewer
func Flash(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	id := strings.Split(url[len(url)-1], ".")[0]
	o, e := storage.GetMimeType(id)
	if e != nil || o.MimeType != "application/x-shockwave-flash" {
		redirect(w, r)
		return
	}
	f, e := os.Open("./assets/fv.html")
	if e != nil {
		redirect(w, r)
		log.Println(errors.Annotate(e, "FlashViewer Open"))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, e = io.Copy(io.Writer(w), io.Reader(f))
	if e != nil {
		redirect(w, r)
		log.Println(errors.Annotate(e, "FlashViewer Copy"))
		return
	}
	w.WriteHeader(200)
	return
}
