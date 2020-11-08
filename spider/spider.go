package spider

import (
	"fmt"
	"github.com/axgle/mahonia"
	"go_spider/headless"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func Download_img(url string, dir_path string, i int, referer string) {
	name := fmt.Sprintf("%s/%d.jpg", dir_path, i)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("if-modified-since", "Mon, 18 May 2020 12:58:04 GMT")
	req.Header.Set("if-none-match", "5ec2865c-32d96")
	req.Header.Set("referer", "http://yishesp.com/xurl_881393.html")
	req.Header.Set("upgrade-insecure-requests", "1")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("img get fail")
		return
	}
	f, err := os.Create(name)
	if err != nil {
		fmt.Println("file creat failed")
		return
	}
	io.Copy(f, res.Body)
	fmt.Println(url + "\t下载成功")
}
func generateURL(rootURL string, i int) (url string) {
	if i == 1 {
		url = fmt.Sprintf("%s.html", rootURL)
	} else if i > 1 {
		url = fmt.Sprintf("%s_%d.html", rootURL, i)
	}
	return url
}
func getData(data string, label string) string {
	start := fmt.Sprintf("<%s>", label)
	stop := fmt.Sprintf("</%s>", label)
	tmp := strings.Replace(data, start, "", 1)
	fin := strings.Replace(tmp, stop, "", 1)
	return fin
}
func getSrc(str string) string {
	len := strings.Index(str, ".jpg")
	if len <= 1 {
		return ""
	}
	len += 4
	ret := string([]byte(str)[:len])
	return ret
}
func downloadAllImages(url string, path string) {
	fmt.Println("访问：", url, "保存到：", path)
	html := headless.GetHTMLByChromebp(url, true) // 通过chrome 获取html

	/* 创建文件夹 */
	matchTitle := regexp.MustCompile(`<title>.*?</title>`)
	matchedTitle_str := matchTitle.FindAllString(html, -1)
	dirName := getData(matchedTitle_str[0], "title")
	dirPATH := filepath.Join(path, dirName)
	_, err := os.Stat(dirPATH)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("exist")
			fmt.Println(err)
			return
		}
		fmt.Println("no exist")
	}
	os.MkdirAll(dirPATH, os.ModePerm)
	/* 匹配图片 */
	match := regexp.MustCompile(`<img src="(.*?)" data-bd-imgshare-binded="1">`)
	matched_str := match.FindAllString(html, -1)
	i := 0
	for _, match_str := range matched_str {
		fmt.Println(dirPATH, "    ", match_str)
		i++
		name := match.FindStringSubmatch(match_str)[1]
		src := getSrc(name)
		time.Sleep(time.Second)
		Download_img(src, dirPATH, i, url)
	}
}
func SpiderImage(rootURL string, i int, dirRoot string, web string) {
	if i < 1 {
		fmt.Println("err page %d", i)
		return
	}
	url := generateURL(rootURL, i)
	fmt.Println("Enter URL: ", url);
	response, err1 := http.Get(url)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	content, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println(err1)
		return
	}
	defer response.Body.Close()
	html := string(content)
	/* 创建文件夹 */
	//html = ConvertToString(html, "gbk", "utf-8")
	matchTitle := regexp.MustCompile(`<title>.*?</title>`)
	matchedTitle_str := matchTitle.FindAllString(html, -1)
	dirName := getData(matchedTitle_str[0], "title")
	dirPATH := filepath.Join(dirRoot, dirName) //创建页目录
	os.MkdirAll(dirPATH, os.ModePerm)
	/* 匹配超链接 */
	match := regexp.MustCompile(`<a href="(.*?)" target="_blank">`)
	matched_str := match.FindAllString(html, -1)
	fmt.Println("Get list OK")
	for _, match_str := range matched_str {
		/* 获取每一个有用的超链接 */
		href := match.FindStringSubmatch(match_str)[1]
		if !strings.Contains(href, "/xurl") {
			continue
		}
		fullhref := fmt.Sprintf("%s%s", web, href) // 拼接url
		/* 访问超链接 */
		downloadAllImages(fullhref, dirPATH)
	}
}

/**************************dowm load txt*******************************/
func getTxtByHtml(html string) string {
	start := strings.Index(html, `<div class="pics"><div id="pic_text_top"></div>`)
	stop := strings.Index(html, `<div id="pic_text_bottom"></div></div>`)
	if start < 0 || stop < 0 {
		return ""
	}
	data := html[start:stop]
	i := strings.Index(data, `</iframe></div></span>`)
	j := strings.Index(data, `<span id="span_ed8">`)
	if i < 0 || j < 0 {
		return ""
	}
	return strings.Replace(data[i+23:j], "<br>", "\n", -1)
}
func downloadTxt(url string, path string) {
	fmt.Println("访问：", url, "保存到：", path)
	html := headless.GetHTMLByChromebp(url, true) // 通过chrome 获取html
	/* 创建文件 */
	matchTitle := regexp.MustCompile(`<title>.*?</title>`)
	matchedTitle_str := matchTitle.FindAllString(html, -1)
	txtName := getData(matchedTitle_str[0], "title")
	/* 修改文件名称 */
	txtName = strings.Replace(txtName, "【", "", -1)
	txtName = strings.Replace(txtName, "】", "", -1)
	fileName := path + "/" + strings.Replace(txtName, " yishesp.com", "", -1) + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("create file failed: ", err)
		file.Close()
		return
	}
	fmt.Println("get txt name:", fileName)
	/* 获取txt */
	txt := getTxtByHtml(html)
	/* 保存txt */
	file.Write([]byte(txt))
	file.Close()
}
	/**
	 * rootURL : the prefix of url
	 * i: the id of html
	 * dirRoot : the save path
	 */
func SpiderTxt(rootURL string, i int, dirRoot string, web string) {
	if i <= 1 {
		fmt.Println("err page %d", i)
		return
	}
	url := generateURL(rootURL, i)
	fmt.Println("Enter URL: ", url);
	response, err1 := http.Get(url)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	content, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println(err1)
		return
	}
	defer response.Body.Close()
	html := string(content)
	/* 创建文件夹 */
	//html = ConvertToString(html, "gbk", "utf-8")
	matchTitle := regexp.MustCompile(`<title>.*?</title>`)
	matchedTitle_str := matchTitle.FindAllString(html, -1)
	dirName := getData(matchedTitle_str[0], "title")
	dirPATH := filepath.Join(dirRoot, dirName) //创建页目录
	os.MkdirAll(dirPATH, os.ModePerm)
	/* 匹配超链接 */
	match := regexp.MustCompile(`<a href="(.*?)" target="_blank">`)
	matched_str := match.FindAllString(html, -1)
	for _, match_str := range matched_str {
		/* 获取每一个有用的超链接 */
		href := match.FindStringSubmatch(match_str)[1]
		if !strings.Contains(href, "/xurl") {
			continue
		}
		fullhref := fmt.Sprintf("%s%s", web, href)
		fmt.Println("Get TXT URL", fullhref)
		/* get href html by chrome */
		downloadTxt(fullhref, dirPATH)
	}
}