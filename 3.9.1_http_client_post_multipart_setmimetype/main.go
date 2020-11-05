package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "hello!")

	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "image/jpeg")
	part.Set("Content-Disposition", `form-data; name="file"; filename="sample.jpg"`)

	fileWriter, err := writer.CreatePart(part)
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
