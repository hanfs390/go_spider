package main
 
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/axgle/mahonia"
)
 
var workResultLock sync.WaitGroup
 
func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
 
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
 
func download_img(request_url string, name string, dir_path string) {
	image, err := http.Get(request_url)
	check(err)
	image_byte, err := ioutil.ReadAll(image.Body)
	defer image.Body.Close()
	//file_path := filepath.Join(dir_path, name+".jpg")
	file_path := filepath.Join(dir_path, "1.jpg")
	fmt.Println("file ", file_path)
	err = ioutil.WriteFile(file_path, image_byte, 0644)
	check(err)
	fmt.Println(request_url + "\t下载成功")
}
 
func spider(i int, dir_path string) {
	defer workResultLock.Done()
	url := "http://pic.netbian.com/4kdongman/index.html"
	response, err2 := http.Get(url)
	check(err2)
	content, err3 := ioutil.ReadAll(response.Body)
	check(err3)
	defer response.Body.Close()
	html := string(content)
	html = ConvertToString(html, "gbk", "utf-8")
	fmt.Println(html)
	match := regexp.MustCompile(`<img src="(.*?)" alt="(.*?)">`)
	matched_str := match.FindAllString(html, -1)
	for _, match_str := range matched_str {
		fmt.Println(match_str)
		var img_url string
		name := match.FindStringSubmatch(match_str)[1]
		src := match.FindStringSubmatch(match_str)[2]
		if strings.HasPrefix(src, "http") != true {
			var buffer bytes.Buffer
			buffer.WriteString("http://pic.netbian.com")
			buffer.WriteString(name)
			img_url = buffer.String()
		} else {																																								            img_url = src
		}
		download_img(img_url, name, dir_path)
	}
}
 
func main() {
	start := time.Now()
	dir := filepath.Dir(".")
	dir_path := filepath.Join(dir, "images")
	err1 := os.MkdirAll(dir_path, os.ModePerm)
	check(err1)

	workResultLock.Add(1)
	go spider(0, dir_path)

	workResultLock.Wait()
	fmt.Println(time.Now().Sub(start))
}
