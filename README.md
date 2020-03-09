# images
the images bed repo  

## 图床使用而已
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/imgbed4github"
	"io/ioutil"
	"log"
)

var token = "token string"

func main() {
	gh := imgbed4github.NewGithub(token, "imgbed", "images")
	res,err := gh.UploadFromLocalFile("/path/to/img.jpg","img")
	if err!=nil {
		log.Println(err.Error())
		return
	}
	var resp imgbed4github.RepositoryContentResponse
	res.Bind(&resp)

	fmt.Println(*resp.Content.Name, *resp.Content.HTMLURL)
}
```