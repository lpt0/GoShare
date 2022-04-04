package handlers

import (
	"fmt"
	"github.com/lpt0/goshare/config"
	"github.com/lpt0/goshare/storage"
	"io"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/juju/errors"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randomName generates a random file name, and checks it.
// If the file exists, it calls itself until it finds one that doesn't, returning that.
func randomName() string {
	var n string
	for i := 0; i < 6; i++ {
		n += string(alphabet[rand.Intn(52)])
	}
	// TODO: Figure out how to properly pass a DB
	if storage.IDExists(n) {
		return randomName()
	}
	log.Printf("Random ID generated: %s\n", n)
	return n
}

func fileUpload(upload multipart.File, h *multipart.FileHeader, e error, uploader string) (string, error) {
	if e != nil {
		return "", e
	}
	var name string
	var mtype string
	id := randomName()
	ext := strings.Split(h.Filename, ".")
	// Check whether it has an extension (split on dot will be bigger than 1)
	if len(ext) > 1 {
		name = id + "." + ext[len(ext)-1]
		mtype = mime.TypeByExtension("." + ext[len(ext)-1])
		log.Printf("Name %s, mtype %s\n", name, mtype)
	} else {
		name = id
		mtype = "application/octet-stream" // Default content-type
	}
	file, e := os.Create(config.FilePath + name)
	if e != nil {
		return "", e
	}
	defer file.Close()
	_, e = io.Copy(io.Writer(file), io.Reader(upload))
	if e != nil {
		return "", e
	}
	r, e := storage.AddUpload(storage.Object{ID: id, Type: storage.File, OriginalName: h.Filename, Location: file.Name(), MimeType: mtype, Uploader: uploader})
	log.Printf("Result: %v, Error: %v\n", r, e)
	if mtype == "application/x-shockwave-flash" {
		return "flash/" + name, nil
	}
	return name, nil
}

func shortenURL(url, uploader string) (string, error) {
	id := randomName()
	// Execute wget with warc here
	c := exec.Command("wget", "--page-requisites", "--delete-after", "--no-directories", "--warc-file="+id, "--warc-cdx", url)
	c.Dir = config.FilePath
	o, e := c.CombinedOutput()
	log.Printf("shortenURL wget: %v\n", string(o))
	if e != nil {
		return "", e
	}
	r, e := storage.AddUpload(storage.Object{ID: id, Type: storage.URL, Location: url, Uploader: uploader})
	log.Printf("Result: %v, Error: %v\n", r, e)
	return id, e
}

// Upload is the file upload handler.
func Upload(w http.ResponseWriter, r *http.Request) {
	var e error
	var name string
	uploader := config.Authorization[strings.ToLower(r.Header.Get("Authorization"))]
	if uploader != "" {
		if r.Header.Get("Content-Type") != "" {
			if strings.Split(r.Header.Get("Content-Type"), ";")[0] == "multipart/form-data" {
				u, h, e := r.FormFile("file")
				name, e = fileUpload(u, h, e, uploader)
			}
		}
		if name == "" && r.Header.Get("Location") != "" {
			name, e = shortenURL(r.Header.Get("Location"), uploader)
		}
		if e != nil {
			w.WriteHeader(500)
			log.Println(errors.Annotate(e, "Upload error"))
			return
		}
		url := config.Protocol + "://" + config.Host
		if config.Port != "" {
			url += ":" + config.Port
		}
		url += "/" + name
		log.Println(uploader + " uploaded " + url)
		fmt.Fprintf(w, url)
		return
	}
	w.WriteHeader(403)
	return
}
