package jsnm

import (
	"fmt"
	"strconv"
)

type RawData struct {
	raw interface{}
}

func NewRawData(raw interface{}) RawData {
	return RawData{raw: raw}
}

func (d *RawData) String() string {
	return fmt.Sprintf("%v", d.raw)
}

func (d *RawData) Int64() (int64, error) {
	return strconv.ParseInt(d.String(), 10, 0)
}

func (d *RawData) Float64() (float64, error) {
	return strconv.ParseFloat(d.String(), 0)
}
