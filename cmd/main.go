package main

import (
	"encoding/json"
	"fmt"
	imgbed4github "github.com/gohouse/imgbed4github"
	"io/ioutil"
	"log"
)

var token string
var gh *imgbed4github.Github

func init() {
	b, _ := ioutil.ReadFile("config.json")
	var tmp struct {
		Token string
	}
	json.Unmarshal(b, &tmp)
	token = tmp.Token
	gh = imgbed4github.NewGithub(token, "imgbed", "images")
}

func main() {
	res,err := gh.UploadFromLocalFile("~/Downloads/img/ie_joke.png","")
	if err!=nil {
		log.Println(err.Error())
		return
	}
	var resp imgbed4github.RepositoryContentResponse
	res.Bind(&resp)

	fmt.Println(*resp.Content.Name, *resp.Content.HTMLURL)
}
func getcontet() {
	var f = "README.md"
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err.Error())
	}
	res, err := gh.Upload(f, &b)
	if err != nil {
		log.Fatal(err.Error())
	}

	var resp imgbed4github.RepositoryContentResponse
	res.Bind(&resp)
	fmt.Println(*resp.Content.Path)
	fmt.Println(*resp.Content.HTMLURL)
}
