package handlers

import (
	"goshare/config"
	"goshare/storage"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/juju/errors"
)

// Default provides the default route (/)
func Default(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	id := strings.Split(url[len(url)-1], ".")[0]
	upload, e := storage.GetUpload(id)
	if e != nil {
		redirect := config.Redirects[rand.Intn(len(config.Redirects))]
		log.Println("Redirecting to " + redirect)
		w.Header().Set("Location", redirect)
		w.WriteHeader(301)
		return
	}
	if upload.Type == storage.URL {
		w.Header().Set("Location", upload.Location)
		w.WriteHeader(301)
		w.Write(make([]byte, 1))
	} else if upload.Type == storage.File {
		f, e := os.Open(upload.Location)
		if e != nil {
			log.Println(errors.Annotate(e, "DefaultRoute File"))
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", upload.MimeType)
		_, e = io.Copy(io.Writer(w), io.Reader(f))
		if e != nil {
			log.Println(errors.Annotate(e, "DefaultRoute Copy"))
		}
		w.WriteHeader(200)
		return
	}
}
