package jsnm

import (
	"encoding/json"
	"github.com/shaalx/goutils"
	"io"
	"os"
)

func BytesFmt(bs []byte) *Jsnm {
	if nil == bs {
		return nil
	}
	v := NewJsnm(nil)
	err := json.Unmarshal(bs, &v.raw.raw)
	if goutils.CheckErr(err) {
		return nil
	}
	if map_data, ok := v.raw.raw.(map[string]interface{}); ok {
		v.data = map_data
	}
	return v
}

func ReaderFmt(r io.Reader) *Jsnm {
	v := NewJsnm(nil)
	err := json.NewDecoder(r).Decode(&v.raw.raw)
	if goutils.CheckErr(err) {
		return nil
	}
	if map_data, ok := v.raw.raw.(map[string]interface{}); ok {
		v.data = map_data
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
