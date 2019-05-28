package jsnm

import (
	"fmt"
	"reflect"
	"strconv"
)

type RawData struct {
	raw interface{}
}

func NewRawData(raw interface{}) RawData {
	return RawData{raw: raw}
}

func (d *RawData) Raw() interface{} {
	if d == nil {
		return nil
	}
	return d.raw
}

func (d *RawData) String() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("%v", d.raw)
}

func (d *RawData) Int64() (int64, error) {
	return strconv.ParseInt(d.String(), 10, 0)
}

func (d *RawData) MustInt64() int64 {
	i64, err := strconv.ParseInt(d.String(), 10, 0)
	if err != nil {
		return 0
	}
	return i64
}

func (d *RawData) MustFloat64() float64 {
	f, err := d.Float64()
	if err != nil {
		return 0.0
	}
	return f
}

func (d *RawData) Float64() (float64, error) {
	return strconv.ParseFloat(d.String(), 0)
}

func (d *RawData) Decode() string {
	typ := reflect.TypeOf(d.raw)
	switch typ.Kind() {
	case reflect.String:
		return d.String()
	case reflect.Float64, reflect.Float32:
		return strconv.FormatInt(int64(d.MustFloat64()), 10)
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(d.MustInt64(), 10)
	default:
		return d.String()
	}
}
