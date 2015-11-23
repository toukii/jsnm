package jsnm

import (
	"fmt"
)

type RawData struct {
	raw interface{}
}

func NewRawData(raw interface{}) RawData {
	return RawData{raw: raw}
}

func (d RawData) String() string {
	return fmt.Sprintf("{raw_data:%#v}", d)
}
