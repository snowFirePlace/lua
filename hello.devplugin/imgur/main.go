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
	"time"

	"github.com/pkg/browser"
)

const version = "0.0.1"

var im imgur
var (
	f        *os.File
	infoLog  *log.Logger
	errorLog *log.Logger
)

type Request struct {
	Data struct {
		Link       string `json:"link"`
		Deletehash string `json:"deletehash"`
	} `json:"data"`
}

type imgur struct {
	Buf        *bytes.Buffer
	Writer     *multipart.Writer
	URL        string
	Deletehash string
}

func (i *imgur) init() {
	i.Buf = new(bytes.Buffer)
	i.Writer = multipart.NewWriter(i.Buf)
}
func logConf() {

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLog = log.New(f, "", log.Ldate|log.Ltime)
	errorLog = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	if len(os.Args) <= 1 {
		return
	}

	logConf()
	defer f.Close()

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	imagePath := os.Args[1]

	log.SetLevel("Debug")

	image, err := os.Open(imagePath)

	if err != nil {
		log.Fatal(err)
	}
	im.upload(image)
	// url https://yandex.ru/images/touch/search?url=***&rpt=imageview&crop=0%3B0%3B1%%3B1
	// https://yandex.ru/images/touch/search?url=https://i.imgur.com/dPrUtO2.jpg&rpt=imageview&crop=0%3B0%3B1%%3B1
	// https://yandex.ru/images/touch/search?url=https://i.imgur.com/qMD15cz.jpg&rpt=imageview
	browser.OpenURL(fmt.Sprintf("https://yandex.ru/images/touch/search?url=%s&rpt=imageview", im.URL))
	time.Sleep(time.Second * 5)
	im.delete()
}
func (i *imgur) upload(image io.Reader) {
	i.init()

	part, _ := i.Writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)
	i.Writer.Close()

	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", i.Buf)
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", ClientID))
	req.Header.Set("Content-Type", i.Writer.FormDataContentType())

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	r := Request{}
	json.Unmarshal(body, &r)
	fmt.Printf("Operation: %s", r.Data.Link)

	i.URL, i.Deletehash = r.Data.Link, r.Data.Deletehash
}
func (i *imgur) delete() {
	i.init()

	err := i.Writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.imgur.com/3/image/%s", i.Deletehash), i.Buf)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", ClientID))
	req.Header.Set("Content-Type", i.Writer.FormDataContentType())

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
