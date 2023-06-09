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
	"strings"
	"time"

	"github.com/pkg/browser"
)

const version = "0.0.1"

var im imgur
var (
	f                *os.File
	infoLog          *log.Logger
	errorLog         *log.Logger
	criticalErrorLog *log.Logger
)

type Request struct {
	Data struct {
		Link       string `json:"link"`
		Deletehash string `json:"deletehash"`
	} `json:"data"`
	Status int `json:""`
}

type imgur struct {
	Buf        *bytes.Buffer
	Writer     *multipart.Writer
	URL        string
	Deletehash string
}

func (i *imgur) newBuf() {
	i.Buf = new(bytes.Buffer)
	i.Writer = multipart.NewWriter(i.Buf)
}
func logConf() {
	path, err := os.Getwd()
	if err != nil {
		log.Print(err)
	}
	f, err := os.OpenFile(path+string(os.PathSeparator)+"imgur.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLog = log.New(f, "", log.Ldate|log.Ltime)
	errorLog = log.New(f, "ERROR\t", log.Ldate|log.Ltime)
	criticalErrorLog = log.New(f, "CRITICAL ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}
func getNameFile(p string) (n string) {
	n = p[strings.LastIndex(p, string(os.PathSeparator))+1:]
	return n
}
func main() {
	logConf()
	defer f.Close()
	if len(os.Args) <= 1 {
		errorLog.Println("Launching the program without pointing to the image is impossible.")
		return
	}

	imagePath := os.Args[1]
	image, err := os.Open(imagePath)

	if err != nil {
		criticalErrorLog.Println(err)
		return
	}
	stat, _ := image.Stat()
	name := stat.Name()
	im.upload(image)
	if im.URL == "" {
		errorLog.Println("Can't upload image to imgur server")
		return
	}

	browser.OpenURL(fmt.Sprintf("https://yandex.ru/images/touch/search?url=%s&rpt=imageview", im.URL))
	time.Sleep(time.Second * 5)
	if im.delete() != 200 {
		errorLog.Println("File not deleted")
		return
	}
	infoLog.Printf("File %s upload to imgur. File delete.\n", name)

}
func (i *imgur) upload(image io.Reader) {
	i.newBuf()

	part, _ := i.Writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)
	i.Writer.Close()

	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", i.Buf)
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", clientID))
	req.Header.Set("Content-Type", i.Writer.FormDataContentType())

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	r := Request{}
	json.Unmarshal(body, &r)
	i.URL, i.Deletehash = r.Data.Link, r.Data.Deletehash
}
func (i *imgur) delete() int {
	i.newBuf()

	err := i.Writer.Close()
	if err != nil {
		criticalErrorLog.Println(err)
		return 400
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.imgur.com/3/image/%s", i.Deletehash), i.Buf)
	if err != nil {
		criticalErrorLog.Println(err)
		return 400
	}
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", clientID))
	req.Header.Set("Content-Type", i.Writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		criticalErrorLog.Println(err)
		return 400
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		criticalErrorLog.Println(err)
		return 400
	}
	r := struct {
		Data   string `json:"data"`
		Status int    `json:"status"`
	}{}
	json.Unmarshal(body, &r)
	return r.Status
}
