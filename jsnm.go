package jsnm

import (
	"fmt"
)

func init() {
	fmt.Println("main...")
}

// No Cache Get
func (j *Jsnm) NCGet(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	// first step: get data from mapdata
	cur, ok := j.data[path[0]]
	if !ok {
		return nil
	}
	// second step: cache the data
	var will_data *Jsnm
	if v, ok := cur.(map[string]interface{}); ok {
		will_data = NewJsnm(v)
	} else {
		will_data = NewRawJsnm(cur)
	}
	if len(path) == 1 {
		return will_data
	}
	return will_data.NCGet(path[1:]...)
}

func (j *Jsnm) Get(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return j
	}
	// first step: get data from cache
	if cache_data, ok := j.cache[path[0]]; ok {
		// fmt.Printf("******cache %s*****\n", path)
		if len(path) == 1 {
			return cache_data
		} else {
			return cache_data.Get(path[1:]...)
		}
	}
	// second step: get data from mapdata
	cur, ok := j.data[path[0]]
	if !ok {
		return nil
	}
	// third step: cache the data
	var will_cache_data *Jsnm
	if v, ok := cur.(map[string]interface{}); ok {
		will_cache_data = NewJsnm(v)
	} else {
		will_cache_data = NewRawJsnm(cur)
	}
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
	arr, ok := (j.raw.raw).([]interface{})
	if !ok {
		return nil
	}
	ret := make([]*Jsnm, 0, len(arr))
	for _, vry := range arr {
		if map_data, ok := vry.(map[string]interface{}); ok {
			ret = append(ret, NewJsnm(map_data))
		} else {
			ret = append(ret, NewRawJsnm(vry))
		}
	}
	return ret
}

func (j *Jsnm) ArrLoc(i int) *Jsnm {
	if j == nil {
		return nil
	}
	arr, ok := (j.raw.raw).([]interface{})
	if !ok {
		return nil
	}
	if i >= len(arr) {
		return nil
	}
	if map_data, ok := arr[i].(map[string]interface{}); ok {
		return NewJsnm(map_data)
	} else {
		return NewRawJsnm(arr[i])
	}
	return nil
}
