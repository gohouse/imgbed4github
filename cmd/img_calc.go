package main

import (
	"encoding/json"
	"github.com/gohouse/imgbed4github"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func initGh() *imgbed4github.Github {
	b, _ := ioutil.ReadFile("config.json")
	var tmp struct {
		Token string
	}
	json.Unmarshal(b, &tmp)
	token := tmp.Token
	return imgbed4github.NewGithub(token, "imgbed", "images")
}
func main() {
	img2rds()
}
func img2rds() {
	//var rds = BootRedis()
	var gh = initGh()
	resp, err := gh.GetContents("img")
	if err != nil {
		log.Println(err.Error())
		return
	}
	var cont []imgbed4github.RepositoryContent
	resp.Bind(&cont)
	for _, item := range cont {
		if *item.Type == "dir" {
			resp2, err := gh.GetContents(*item.Path)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var cont2 []imgbed4github.RepositoryContent
			resp2.Bind(&cont2)
			for _, item2 := range cont2 {
				//log.Println("item2.Path: ", *item2.Path)
				//log.Printf("url: %s%s\n", gh.UrlPrefix(), *item2.Path)
				//log.Println("sha1: ", strings.TrimRight(*item2.Name, path.Ext(*item2.Name)))

				var sha1str = strings.TrimRight(*item2.Name, path.Ext(*item2.Name))
				var imgurl = gh.UrlPrefix() + *item2.Path
				//// 放入redis
				//_, err := rds.SAdd(SetImg, sha1str).Result()
				//if err != nil {
				//	log.Println("err: ", sha1str, err.Error())
				//}
				//_, err = rds.HSet(HashImg, sha1str, imgurl).Result()
				//if err != nil {
				//	log.Println("err: ", sha1str, err.Error())
				//}
				////return
			}
		}
	}
}
func gif2rds() {
	//var rds = BootRedis()
	var gh = initGh()
	resp, err := gh.GetContents("gif")
	if err != nil {
		log.Println(err.Error())
		return
	}
	var cont []imgbed4github.RepositoryContent
	resp.Bind(&cont)
	for _, item := range cont {
		if *item.Type == "dir" {
			resp2, err := gh.GetContents(*item.Path)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var cont2 []imgbed4github.RepositoryContent
			resp2.Bind(&cont2)
			for _, item2 := range cont2 {
				//log.Println("item2.Path: ", *item2.Path)
				//log.Printf("url: %s%s\n", gh.UrlPrefix(), *item2.Path)
				//log.Println("sha1: ", strings.TrimRight(*item2.Name, path.Ext(*item2.Name)))

				var sha1str = strings.TrimRight(*item2.Name, path.Ext(*item2.Name))
				var imgurl = gh.UrlPrefix() + *item2.Path
				//// 放入redis
				//_, err := rds.SAdd(SetGif, sha1str).Result()
				//if err != nil {
				//	log.Println("err: ", sha1str, err.Error())
				//}
				//_, err = rds.HSet(HashGif, sha1str, imgurl).Result()
				//if err != nil {
				//	log.Println("err: ", sha1str, err.Error())
				//}
				//return
			}
		}
	}
}

//const (
//	SetGif  = "s:gif"
//	SetImg  = "s:img"
//	HashGif = "h:gif"
//	HashImg = "h:img"
//)
//
//func BootRedis() *redis.Client {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "123456",
//	})
//	return client
//}
