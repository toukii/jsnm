package jsnm

import (
	"net/http"
	"fmt"
)

// http://news-at.zhihu.com/api/3/news/hot

func ExampleWriter() {
	resp, err := http.Get("http://news-at.zhihu.com/api/3/news/hot")
	if checkerr(err) {
		fmt.Println(err)
		return
	}
	jm := ReaderFmt(resp.Body)
	// t.Log(jm.RawData())
	arr := jm.Get("recent").Arr()
	fmt.Printf("arr's length:%d\n", len(arr))
	for i, it := range arr {
		fmt.Println(i+1, it.Get("title").RawData().String())
		fmt.Println(i+1, it.Get("url").RawData().String())
	}
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
