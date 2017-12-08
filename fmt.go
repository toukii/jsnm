package jsnm

import (
	"encoding/json"
	"io"
	"os"

	"github.com/everfore/exc"
	"github.com/toukii/goutils"
)

func BytesFmt(bs []byte) *Jsnm {
	if nil == bs {
		return nil
	}
	v := NewJsnm(nil)
	err := json.Unmarshal(bs, &v.raw_data)
	if goutils.CheckNoLogErr(err) {
		return nil
	}
	return v
}

func ReaderFmt(r io.Reader) *Jsnm {
	v := NewJsnm(nil)
	err := json.NewDecoder(r).Decode(&v.raw_data)
	if goutils.CheckNoLogErr(err) {
		return nil
	}
	return v
}

func FileNameFmt(fn string) *Jsnm {
	rf, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	if goutils.CheckNoLogErr(err) {
		return nil
	}
	return ReaderFmt(rf)
}

func CmdFmt(cmd string) *Jsnm {
	bs, err := exc.NewCMD(cmd).DoNoTime()
	if goutils.CheckErr(err) {
		return nil
	}
	return BytesFmt(bs)
}

func StructFmt(bs []byte, st interface{}) error {
	err := json.Unmarshal(bs, st)
	if goutils.CheckErr(err) {
		return err
	}
	return nil
}
