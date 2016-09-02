package jsnm

import (
	"encoding/json"
	"github.com/toukii/goutils"
	"io"
	"os"
)

func BytesFmt(bs []byte) *Jsnm {
	if nil == bs {
		return nil
	}
	v := NewJsnm(nil)
	err := json.Unmarshal(bs, &v.raw_data)
	if goutils.CheckErr(err) {
		return nil
	}
	return v
}

func ReaderFmt(r io.Reader) *Jsnm {
	v := NewJsnm(nil)
	err := json.NewDecoder(r).Decode(&v.raw_data)
	if goutils.CheckErr(err) {
		return nil
	}
	return v
}

func FileNameFmt(fn string) *Jsnm {
	rf, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if goutils.CheckErr(err) {
		return nil
	}
	return ReaderFmt(rf)
}
