package handlers

import (
	"fmt"
	"goshare/storage"
	"io"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Path is the storage path for uploaded files
var path = os.Getenv("FILES") + "/"

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

func fileUpload(upload multipart.File, h *multipart.FileHeader, e error) (string, error) {
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
	file, e := os.Create(path + name)
	if e != nil {
		return "", e
	}
	defer file.Close()
	_, e = io.Copy(io.Writer(file), io.Reader(upload))
	if e != nil {
		return "", e
	}
	r, e := storage.AddUpload(storage.Object{ID: id, Type: storage.File, OriginalName: h.Filename, Location: file.Name(), MimeType: mtype})
	log.Printf("Result: %v, Error: %v\n", r, e)
	return name, nil
}

func shortenURL(url string) (string, error) {
	id := randomName()
	// Execute wget with warc here
	c := exec.Command("wget", "--page-requisites", "--delete-after", "--no-directories", "--warc-file="+id, "--warc-cdx", url)
	c.Dir = path
	o, e := c.CombinedOutput()
	log.Printf("shortenURL wget: %v\n", string(o))
	if e != nil {
		return "", e
	}
	r, e := storage.AddUpload(storage.Object{ID: id, Type: storage.URL, Location: url})
	log.Printf("Result: %v, Error: %v\n", r, e)
	return id, e
}

// Upload is the file upload handler.
func Upload(w http.ResponseWriter, r *http.Request) {
	var e error
	var name string
	if r.Header.Get("Content-Type") != "" {
		if strings.Split(r.Header.Get("Content-Type"), ";")[0] == "multipart/form-data" {
			log.Println(r.Header.Get("Content-Type"))
			name, e = fileUpload(r.FormFile("file"))
		}
	}
	if name == "" && r.Header.Get("Location") != "" {
		name, e = shortenURL(r.Header.Get("Location"))
	}
	if e != nil {
		w.WriteHeader(500)
		log.Printf("Upload error: %v\n", e)
		fmt.Fprintf(w, e.Error())
	}
	url := "http://" + r.Host + "/" + name
	log.Println(url)
	fmt.Fprintf(w, url)
}
