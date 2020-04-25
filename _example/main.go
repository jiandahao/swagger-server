package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"net/http"
)

func main() {
	//client := http.DefaultClient
	file, err := ioutil.ReadFile("test_swagger.swagger.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	formFile, err  := w.CreateFormFile("file", "test_swagger.swagger.json")
	if err != nil {
		panic(err)
	}

	formFile.Write(file)

	contentType := w.FormDataContentType()

	// should be closed before sending request
	w.Close()

	resp, err := http.Post("http://127.0.0.1:8088/new", contentType, buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}