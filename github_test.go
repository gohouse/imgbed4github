package imgbed4github

import (
	"encoding/json"
	"flag"
	"github.com/gohouse/golib/file"
	t2 "github.com/gohouse/golib/t"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// https://developer.github.com/v3/repos/contents/#create-or-update-a-file

var gh *Github
var f string
var token struct{Token string}

func initGh()  {
	flag.StringVar(&f,"f","config.json","配置文件")
	flag.Parse()
	b,err := file.NewFile(f).ReadFile()
	if err!=nil {
		panic(err)
	}
	err = t2.New(b).Bind(&token)
	if err!=nil {
		panic(err)
	}

	gh = NewGithub(token.Token, "imgbed", "images")
}
func TestNewGithub(t *testing.T) {
	initGh()
	t.Log(gh.owner)
}

func TestGithub_Upload(t *testing.T) {
	b, err := ioutil.ReadFile("../../img/2.png")
	if err != nil {
		t.Error(err.Error())
		return
	}

	res,err := gh.Upload("test/a.png", &b)
	var aaa RepositoryContentResponse
	res.Bind(&aaa)
	t.Log(*aaa.Content.Name)
	t.Log(*aaa.Content.Path)
	t.Log(*aaa.Content.HTMLURL)
}

func TestGithub_UploadFromLocalFile(t *testing.T) {
	res,err := gh.UploadFromLocalFile("../../img/2.png","")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(res.String())
}

func TestGithub_UploadFromLocalDir(t *testing.T) {
	var localdir = "../../img"
	res,err := gh.UploadFromLocalDir(localdir,"test2")
	if err != nil {
		t.Error(err.Error())
	}
	var result []*RepositoryContent
	for _,item:=range res{
		var tmp RepositoryContentResponse
		item.Bind(&tmp)
		result = append(result, tmp.Content)
		t.Log(*tmp.Content.Path)
	}
}

func TestGithub_UploadFromLocalDir2(t *testing.T) {
	var localdir = "../../img"
	dirs,err := ioutil.ReadDir(localdir)
	if err!=nil {
		return
	}

	for _,f := range dirs {
		//fullPath := filepath.Join(localdir, f.Name())
		f := filepath.Join("/test2","//", f.Name())
		t.Log(f)
	}
}

func TestGithub_GetContents(t *testing.T) {
	initGh()
	resp,err:=gh.GetContents("gif/3c1/3c1cc42404547205c4ebcc817e38f908de257699.gif")
	if err!=nil {
		t.Error(err.Error())
		return
	}
	var cont []RepositoryContent
	resp.Bind(&cont)
	t.Log(len(cont))
	b,e := json.Marshal(cont)
	t.Log(e)
	t.Logf("%s",b)
}
