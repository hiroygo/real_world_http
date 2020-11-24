package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if rdump, err := httputil.DumpRequest(r, true); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		log.Println(string(rdump))
	}

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	formFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer formFile.Close()

	file, err := os.Create(fileHeader.Filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, formFile); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
