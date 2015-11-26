package jsnm

import (
	"net/http"
	"testing"
)

// http://news-at.zhihu.com/api/3/news/hot

func TestZhihu(t *testing.T) {
	resp, err := http.Get("http://news-at.zhihu.com/api/3/news/hot")
	if checkerr(err) {
		t.Error(err)
	}
	jm := ReaderFmt(resp.Body)
	// t.Log(jm.RawData())
	arr := jm.Get("recent").Arr()
	t.Logf("arr's length:%d\n", len(arr))
	for i, it := range arr {
		t.Log(i+1, it.Get("title").RawData().String())
		t.Log(i+1, it.Get("url").RawData().String())
	}
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
