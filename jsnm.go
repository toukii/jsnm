package jsnm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"reflect"
)

// No Cache Get
func (j *Jsnm) NCGet(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	// first step: get data from mapdata
	//if j.map_data == nil {
	//	j.map_data = make(MapData)
	//	fmt.Println("make map_data")
	map_data, ok := j.raw_data.(map[string]interface{});
	if !ok {
			//j.map_data = map_data
			//fmt.Println("cache map_data")
		//} else {
		return nil
		//}
	}
	cur, ok := map_data[path[0]]
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
			//fmt.Printf("**make cache,subpath:%s, path:%#v\n",subpath,path)
		} else {
			if cache_jm, exist = jm.cache[subpath]; exist {
				jm = cache_jm
				continue
			}
		}
		if jm.map_data == nil {
			jm.map_data = make(MapData)
			if map_data, ok := jm.raw_data.(map[string]interface{}); ok {
				//fmt.Printf("@@cache map_data,subpath:%s, path:%#v\n",subpath,path)
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
		//fmt.Printf("##cache jsnm, subpath:%s\n",subpath)
		jm = sub_jm
	}
	return jm
}

func (j *Jsnm) Get(path string) *Jsnm {
	if j==nil {
		return nil
	}
	// first step: get data from cache
	if nil == j.cache {
		j.cache = make(map[string]*Jsnm)
		//fmt.Println("**make cache")
	} else {
		if cache_data, ok := j.cache[path]; ok {
			return cache_data
		}
	}
	// second step: get data from mapdata
	if j.map_data == nil {
		if map_data, ok := j.raw_data.(map[string]interface{}); ok {
			j.map_data = map_data
			//fmt.Println("@@cache map_data, path:",path)
		} else {
			return nil
		}
	}
	cur, ok := j.map_data[path]
	if !ok {
		return nil
	}
	// third step: cache the data
	will_cache_data := NewJsnm(cur)
	//fmt.Println("##cache jsnm, path:",path)
	j.cache[path] = will_cache_data
	return will_cache_data
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
	//fmt.Println("cache arr_data")
	j.arr_data = ret
	return ret
}


func (j *Jsnm) ArrLocs(locs ...int) *Jsnm {
	if len(locs)<=0 {
		return nil
	}
	subarr:= j.ArrLoc(locs[0])
	l:=len(locs)
	for i := 1; i<l; i++ {
		if subarr==nil {
			 return nil
		}
		subarr = subarr.ArrLoc(locs[i])
	}
	return subarr
}


func (j *Jsnm) ArrLoc(i int) *Jsnm {
	if j == nil {
		return nil
	}
	if nil != j.arr_data {
		if i >= len(j.arr_data) {
			return nil
		}
		if j.arr_data[i] != nil {
			return j.arr_data[i]
		}
	}
	arr, ok := (j.raw_data).([]interface{})
	if !ok {
		return nil
	}
	arr_cache := make([]*Jsnm, len(arr))
	//for _, vry := range arr {
	//	arr_cache = append(arr_cache, NewJsnm(vry))
	//}
	//fmt.Println("cache arr_data")
	j.arr_data = arr_cache
	j.arr_data[i] = NewJsnm(arr[i])
	if i >= len(arr) {
		return nil
	}
	//return arr_cache[i]
	return j.arr_data[i]
}

func (j *Jsnm) ArrPath(path ...interface{}) *Jsnm {
	if len(path) <= 0 {
		return j
	}
	// fmt.Println(path)
	switch reflect.TypeOf(path[0]).Kind() {
	case reflect.String:
		return j.Get(path[0].(string)).ArrPath(path[1:]...)
	case reflect.Int:
		return j.ArrLoc(path[0].(int)).ArrPath(path[1:]...)
	}
	return nil
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

func (j *Jsnm) String() string {
	if nil == j {
		return ""
	}
	return j.RawData().String()
}