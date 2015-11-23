package jsnm

import (
	"encoding/json"
	"fmt"
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

type Jsnm struct {
	data    map[string]interface{}
	subData map[string]*Jsnm
}

func NewJsnm(m map[string]interface{}) *Jsnm {
	return &Jsnm{
		data:    m,
		subData: make(map[string]*Jsnm),
	}
}

func (j *Jsnm) Get(path ...string) (interface{}, error) {
	if len(path) > 0 {
		if sub_cur, ok := j.subData[path[0]]; ok {
			// fmt.Println("******subData", path)
			if len(path) == 1 {
				return sub_cur.data[path[0]], nil
			} else {
				return sub_cur.Get(path[1:]...)
			}
		}
		cur, ok := j.data[path[0]]
		if !ok {
			return nil, fmt.Errorf("ERR:cannot mapping, map[%s]interface{}. ", path[0])
		}
		if len(path) == 1 {
			return cur, nil
		}
		if v, ok := cur.(map[string]interface{}); ok {
			tj := NewJsnm(v)
			j.subData[path[0]] = tj
			return tj.Get(path[1:]...)
		} else {
			return nil, fmt.Errorf("ERR:cannot convert %#v ==> map[string]interface{}.", cur)
		}
	}
	return nil, fmt.Errorf("ERR:path is end.")
}
