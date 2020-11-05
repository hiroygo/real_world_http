package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "hello!")

	fileWriter, err := writer.CreateFormFile("file", "sample.jpg")
	if err != nil {
		panic(err)
	}

	readFile, err := os.Open("sample.jpg")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	io.Copy(fileWriter, readFile)
	writer.Close()

	resp, err := http.Post("http://localhost:8080", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Status:", resp.Status)
}
