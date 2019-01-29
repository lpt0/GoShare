package handlers

import (
	"goshare/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Default provides the default route (/)
func Default(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	id := strings.Split(url[len(url)-1], ".")[0]
	upload, e := storage.GetUpload(id)
	if e != nil {
		w.Header().Set("Location", "https://plaza.one")
		w.WriteHeader(301)
		//w.Write(make([]byte, 1))
	}
	if upload.Type == storage.URL {
		w.Header().Set("Location", upload.Location)
		w.WriteHeader(301)
		w.Write(make([]byte, 1))
	} else if upload.Type == storage.File {
		f, e := os.Open(upload.Location)
		if e != nil {
			log.Printf("DefaultRoute error: %v\n", e)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", upload.MimeType)
		_, e = io.Copy(io.Writer(w), io.Reader(f))
		if e != nil {
			log.Printf("DefaultRouteCopy error: %v\n", e)
		}
		w.WriteHeader(200)
		return
	}
}
