package main

import (
	"fmt"
	"go_spider/spider"
	"os"
	"path/filepath"
	"strconv"
	"time"
)
func downImage(root string, url string, start int, stop int) {
	now := time.Now()
	dir := filepath.Dir(".")
	dir_path := filepath.Join(dir, "images") //创建一级目录
	err := os.MkdirAll(dir_path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := start; i <= stop; i++ {
		spider.SpiderImage(url, i, dir_path, root)
	}
	fmt.Println(time.Since(now))
}
func downTxt(root string, url string, start int, stop int) {
	now := time.Now()
	dir := filepath.Dir(".")
	dir_path := filepath.Join(dir, "txt") //创建一级目录
	err := os.MkdirAll(dir_path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := start; i <= stop; i++ {
		spider.SpiderTxt(url, i, dir_path, root)
	}
	fmt.Println(time.Since(now))
}
func main() {
	//downImage("http://yishesp.com", "http://yishesp.com/piwy_12", 103, 110)
	//downTxt("http://yishesp.com", "http://yishesp.com/piwy_20", 288, 350)
	arg_num := len(os.Args)
	if arg_num != 6 {
		fmt.Printf("the num of input is %d not 6\n",arg_num)
	}
	fmt.Printf("they are :\n")
	for i := 0 ; i < arg_num ;i++{
		fmt.Println(os.Args[i])
	}
	tp,err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("type is err ", err)
		return
	}
	root := os.Args[2]
	hex := os.Args[3]
	start, _ := strconv.Atoi(os.Args[4])
	stop, _ := strconv.Atoi(os.Args[5])
	if tp == 1 {
		downImage(root, hex, start, stop)
	} else if tp == 2 {
		downTxt(root, hex, start, stop)
	}
}

