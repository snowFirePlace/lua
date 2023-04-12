package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

var conf configuration

type configuration struct {
	Client string `yaml:"clientID"`
	Secret string `yaml:"secret"`
}
type Request struct {
	Data struct {
		Link       string `json:"link"`
		Deletehash string `json:"deletehash"`
	} `json:"data"`
}

func main() {
	yfile, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yfile, &conf)

	if err != nil {

		log.Fatal(err)
	}

	// url := "https://api.imgur.com/3/image"
	// method := "POST"
	image, err := os.Open("./imgs/9e63c8436ff97c485da4c3b93d84cb62.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	upload(image)

}
func upload(image io.Reader) (string, string) {
	r := Request{}
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, _ := writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)

	writer.Close()

	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", conf.Client))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &r)
	fmt.Printf("Operation: %s", r.Data.Link)
	return r.Data.Link, r.Data.Deletehash
}
func delete(d string) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.imgur.com/3/image/%s", d), payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", conf.Client))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
