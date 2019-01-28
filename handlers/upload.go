package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
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
	// Query database for it
	return n
}

func fileUpload(u multipart.File, h *multipart.FileHeader, e error) (string, error) {
	n := randomName()
	if e != nil {
		return "", e
	}
	f, e := os.Create("files/" + h.Filename)
	if e != nil {
		return "", e
	}
	defer f.Close()
	_, e = io.Copy(io.Writer(f), io.Reader(u))
	if e != nil {
		return "", e
	}
	log.Println(n)
	return n, nil
}

func shortenURL() (string, error) {
	return "", nil
}

// Upload is the file upload handler.
func Upload(w http.ResponseWriter, r *http.Request) {
	var e error
	var n string
	if r.Header["Content-Type"] == nil {
		w.WriteHeader(500)
		log.Println("Missing Content-Type!")
		fmt.Fprintf(w, "Missing content-type")
	}
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] == "multipart/form-data" {
		log.Println("It's a file")
		n, e = fileUpload(r.FormFile("file"))
	}
	if e != nil {
		w.WriteHeader(500)
		log.Println(e)
		fmt.Fprintf(w, e.Error())
	}
	log.Println("http://" + r.Host + "/" + n + "\n")
	fmt.Fprintf(w, "http://"+r.Host+"/"+n+"\n")
}
