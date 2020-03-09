package imgbed4github

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/golib/curl"
	"github.com/gohouse/imgbed4github/util"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type CommitAuthor struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
type RepositoryContentFileOptions struct {
	Message   *string       `json:"message,omitempty"`
	Content   *string       `json:"content,omitempty"` // unencoded
	SHA       *string       `json:"sha,omitempty"`
	Branch    *string       `json:"branch,omitempty"`
	Author    *CommitAuthor `json:"author,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
}

// RepositoryContent represents a file or directory in a github repository.
type RepositoryContent struct {
	Type *string `json:"type,omitempty"`
	// Target is only set if the type is "symlink" and the target is not a normal file.
	// If Target is set, Path will be the symlink path.
	Target   *string `json:"target,omitempty"`
	Encoding *string `json:"encoding,omitempty"`
	Size     *int    `json:"size,omitempty"`
	Name     *string `json:"name,omitempty"`
	Path     *string `json:"path,omitempty"`
	// Content contains the actual file content, which may be encoded.
	// Callers should call GetContent which will decode the content if
	// necessary.
	Content     *string `json:"content,omitempty"`
	SHA         *string `json:"sha,omitempty"`
	URL         *string `json:"url,omitempty"`
	GitURL      *string `json:"git_url,omitempty"`
	HTMLURL     *string `json:"html_url,omitempty"`
	DownloadURL *string `json:"download_url,omitempty"`
}
type RepositoryContentResponse struct {
	Content *RepositoryContent
}

var host = "https://api.github.com"

type Github struct {
	c                 *curl.Curl
	url               string
	owner, repo, path string
}

func NewGithub(token, owner, repo string) *Github {
	c := curl.NewCurl(curl.ParamHeader(curl.H{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("token %s", token),
	}))
	//host := "https://api.github.com/users/imgbed"
	return &Github{
		c: c,
		//token: token,
		owner: owner,
		repo:  repo,
	}
}

// Upload
// PUT /repos/:owner/:repo/contents/:path
func (gh *Github) Upload(remoteFilename string, content *[]byte) (res *curl.Response, err error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", host, gh.owner, gh.repo, remoteFilename)

	var msg = "img upload"
	var contentbase64 = string(util.Bytes2Base64(*content))

	var rco = RepositoryContentFileOptions{
		Content: &contentbase64,
		Message: &msg,
		Committer: &CommitAuthor{
			Name:  "imgbed",
			Email: "kevin2019010203@gmail.com",
		},
	}
	//sha := gh.GetSHA(remoteFilename)
	//if sha != "" {
	//	return
	//	rco.SHA = &sha
	//}
	b, err := json.Marshal(&rco)
	if err != nil {
		return
	}

	res, err = gh.c.Put(url, curl.ParamJson(&b))
	return
}

// GetContents
// GET /repos/:owner/:repo/contents/:path
func (gh *Github) GetContents(path string) (resp *curl.Response, err error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", host, gh.owner, gh.repo, path)
	return gh.c.Get(url)
}

func (gh *Github) UploadFromUrl(url, remoteDir, ext string, pts ...curl.ParamHandleFunc) (res *curl.Response, err error) {
	resp, err := gh.c.Get(url, pts...)
	if err != nil {
		return
	}
	var b = resp.Bytes()
	filename := BuildFilename(remoteDir, ext, &b)
	return gh.Upload(filename, &b)
}

func (gh *Github) UploadFromLocalFile(localfile, remoteDir string) (res *curl.Response, err error) {
	b, err := ioutil.ReadFile(localfile)
	if err != nil {
		return
	}
	filename := BuildFilename(remoteDir, path.Ext(localfile), &b)
	return gh.Upload(filename, &b)
}

func (gh *Github) UploadFromLocalDir(localdir, remoteDir string) (res []*curl.Response, errs error) {
	dirs, err := ioutil.ReadDir(localdir)
	if err != nil {
		return
	}

	for _, f := range dirs {
		fullPath := filepath.Join(localdir, f.Name())
		resp, err := gh.UploadFromLocalFile(fullPath, remoteDir)
		if err != nil {
			err = fmt.Errorf("Failed to read markdown file %s error: %s ", fullPath, err)
			continue
		} else {
			res = append(res, resp)
		}
	}
	return
}

//func (gh *Github) UploadFromBytes(b []byte,path string) (res *curl.Response ,err error) {
//	str := string(util.Bytes2Base64(b))
//	return gh.Upload(path, &str)
//}

func (gh *Github) GetSHA(path string) string {
	res, err := gh.GetContents(path)

	var sha struct {
		SHA string `json:"sha"`
	}
	if err != nil {
		return ""
	}
	err = json.Unmarshal(res.Bytes(), &sha)
	if err != nil {
		return ""
	}

	return sha.SHA
}

func (gh *Github) UrlPrefix() string {
	//return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/", gh.owner, gh.repo)
	return fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s@master/", gh.owner, gh.repo)
}

func BuildFilename(dir, ext string, content *[]byte) string {
	sha1string := util.Sha1Bytes(content)
	if len(sha1string) < 3 {
		return ""
	}
	// 自动识别是否是图片类型
	exttmp := http.DetectContentType(*content)
	//fmt.Println(exttmp)
	if util.StartWith(exttmp, "image/") {
		arr := strings.Split(exttmp, "/")
		ext = arr[len(arr)-1]
		if dir == "" {
			if ext == "gif" {
				dir = "gif"
			} else {
				dir = "img"
			}
		}
	} else {
		ext = strings.TrimLeft(ext, ".")
		if dir == "" {
			dir = "other"
		}
	}
	return strings.Replace(
		strings.TrimLeft(fmt.Sprintf("%s/%s/%s.%s", dir, sha1string[:3], sha1string, ext), "/"),
		"..", ".", -1)
}
