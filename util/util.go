package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func File2Base64(filearg string) (br []byte, err error) {
	var b []byte
	if b,err = ioutil.ReadFile(filearg);err!=nil {
		return
	}
	br = Bytes2Base64(b)
	return
}

func Bytes2Base64(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}

func Img2Base64(img string) (b64 string, err error) {
	var enc []byte
	if enc,err = File2Base64(img);err!=nil {
		return
	}
	mime := http.DetectContentType(enc)

	if StartWith(mime,"image/") {
		b64 = fmt.Sprintf("data:%s;base64,%s", mime, enc)
	}
	return
}

func StartWith(src,short string) bool {
	if len(src)< len(short) {
		return false
	}

	if src[:len(short)]==short{
		return true
	}
	return false
}

func Md5File(filename string) string {
	f,_ :=  os.Open(filename)
	h := md5.New()
	io.Copy(h,f)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5IoReader(f io.Reader) string {
	h := md5.New()
	io.Copy(h,f)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1File(filename string) string {
	f,_ :=  os.Open(filename)
	h := sha1.New()
	io.Copy(h,f)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1Bytes(b *[]byte) string {
	h := sha1.New()
	io.Copy(h,bytes.NewBuffer(*b))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1String(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
