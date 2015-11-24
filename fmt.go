package jsnm

import (
	"encoding/json"
	"github.com/shaalx/goutils"
	"io"
	"io/ioutil"
	"os"
)

func BytesFmt(bs []byte) *Jsnm {
	if nil == bs {
		return nil
	}
	v := NewJsnm(nil)
	err := json.Unmarshal(bs, &v.data)
	if goutils.CheckErr(err) {
		return nil
	}
	v.raw.raw = v.data
	return v
}

func ReaderFmt(r io.Reader) *Jsnm {
	bs, err := ioutil.ReadAll(r)
	if goutils.CheckErr(err) {
		return nil
	}
	return BytesFmt(bs)
}

func FileNameFmt(fn string) *Jsnm {
	rf, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if goutils.CheckErr(err) {
		return nil
	}
	return ReaderFmt(rf)
}
