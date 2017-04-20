package jsnm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// No Cache Get
func (j *Jsnm) NCGet(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	// first step: get data from mapdata
	if j.map_data == nil {
		j.map_data = make(MapData)
		if map_data, ok := j.raw_data.(map[string]interface{}); ok {
			j.map_data = map_data
		} else {
			return nil
		}
	}
	cur, ok := j.map_data[path[0]]
	if !ok {
		return nil
	}
	// second step: cache the data
	will_data := NewJsnm(cur)
	if len(path) == 1 {
		return will_data
	}
	return will_data.NCGet(path[1:]...)
}

// No Cache Get
func (j *Jsnm) PathGet(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	jm := j
	var cache_jm *Jsnm
	var exist bool
	var sub_data interface{}
	for _, subpath := range path {
		if jm.cache == nil {
			jm.cache = make(map[string]*Jsnm)
		} else {
			if cache_jm, exist = jm.cache[subpath]; exist {
				jm = cache_jm
				continue
			}
		}
		if jm.map_data == nil {
			jm.map_data = make(MapData)
			if map_data, ok := jm.raw_data.(map[string]interface{}); ok {
				jm.map_data = map_data
			} else {
				return nil
			}
		}
		sub_data, exist = jm.map_data[subpath]
		if !exist {
			return nil
		}
		sub_jm := NewJsnm(sub_data)
		jm.cache[subpath] = sub_jm
		jm = sub_jm
	}
	return jm
}

func (j *Jsnm) Get(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	// first step: get data from cache
	if nil == j.cache {
		j.cache = make(map[string]*Jsnm)
	} else {
		if cache_data, ok := j.cache[path[0]]; ok {
			if len(path) == 1 {
				return cache_data
			} else {
				return cache_data.Get(path[1:]...)
			}
		}
	}
	// second step: get data from mapdata
	if j.map_data == nil {
		j.map_data = make(MapData)
		if map_data, ok := j.raw_data.(map[string]interface{}); ok {
			j.map_data = map_data
		} else {
			return nil
		}
	}
	cur, ok := j.map_data[path[0]]
	if !ok {
		return nil
	}
	// third step: cache the data
	will_cache_data := NewJsnm(cur)
	j.cache[path[0]] = will_cache_data
	if len(path) == 1 {
		return will_cache_data
	}
	return will_cache_data.Get(path[1:]...)
}

func (j *Jsnm) Arr() []*Jsnm {
	if j == nil {
		return nil
	}
	if j.arr_data != nil {
		return j.arr_data
	}
	arr, ok := (j.raw_data).([]interface{})
	if !ok {
		return nil
	}
	ret := make([]*Jsnm, 0, len(arr))
	for _, vry := range arr {
		ret = append(ret, NewJsnm(vry))
	}
	j.arr_data = ret
	return ret
}

func (j *Jsnm) ArrLoc(i int) *Jsnm {
	if j == nil {
		return nil
	}
	if nil != j.arr_data {
		if i >= len(j.arr_data) {
			return nil
		}
		return j.arr_data[i]
	}
	arr, ok := (j.raw_data).([]interface{})
	if !ok {
		return nil
	}
	if i >= len(arr) {
		return nil
	}
	return NewJsnm(arr[i])
}

func (j *Jsnm) ArrGet(path ...string) *Jsnm {
	if len(path) <= 0 {
		return j
	}
	// fmt.Println(path)
	if strings.HasSuffix(path[0], "\"") {
		path_0 := strings.Trim(path[0], "\"")
		return j.Get(path_0).ArrGet(path[1:]...)
	}
	if matched, err := regexp.MatchString("\\d", path[0]); matched {
		if err != nil {
			fmt.Println("path-0", err)
		}
		loc, err2 := strconv.ParseInt(path[0], 10, 64)
		if err2 != nil {
			fmt.Println(err2)
		}
		return j.ArrLoc(int(loc)).ArrGet(path[1:]...)
	}
	return j.Get(path[0]).ArrGet(path[1:]...)

}
