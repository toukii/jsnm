package jsnm

import (
	"fmt"
)

func init() {
	fmt.Println("main...")
}

func (j *Jsnm) Get(path ...string) *Jsnm {
	if j == nil || len(path) <= 0 {
		return nil
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
