#	[jsnm][1]

---------------------


__json mapping for map[string]interface{}__

提供缓存的json解析器，初衷是提高重复解析json速度。


## 用法

方法有：__Get PathGet Arr ArrLoc ArrLocs ArrPath__，为了对比，提供了不带缓存的方法：__NCGet__.

 - 带有数组的，要用ArrPath;

 - 无数组形式可以用Get或PathGet;

 - 若不使用缓存，用NCGet,性能会低。


1. 定义结构体如下：

```go
type User struct {
	Name    string
	Age     byte
	Loc     []string
	Friends map[string]*User
}
```

2. 数据如下：

```json
{
	"Name": "foo",
	"Age": 21,
	"Friends": {
		"bar": {
			"Name": "bar",
			"Age": 22,
			"Friends": {
				"bar": null
			},
			"Loc": [
				"shanghai",
				"jiaxing"
			]
		},
		"kaa": null
	},
	"Loc": [
		"beijing",
		"tianjin"
	]
}
```

3. 用法

```
js.Get("Friends").Get("bar").Get("Loc").ArrLoc(0).String()
js.PathGet("Friends","bar").Get("Loc").Arr[0].String()
js.Get("Friends").PathGet("bar","Loc").ArrLoc(0).String()
js.PathGet("Friends","bar","Loc").Arr[0].String()
js.ArrPath("Friends","bar","Loc",0).String()
```

再如：

```
[
    [
        [
            [
                "a"
            ]
        ]
    ]
]
```

获取`a`的方式：

```
js.ArrLocs(0,0,0,0).String()
js.ArrLocs(0,0).ArrLoc(0).Arr[0].String()
js.Arr[0].ArrLocs(0,0).ArrLoc(0).String()
js.ArrLoc(0).Arr[0].ArrLocs(0,0).String()
```

##	数据结构设计

```go
type Jsnm struct {
	raw   RawData
	data  MapData
	arr_data []*Jsnm
	cache map[string]*Jsnm
}

type RawData struct {
	raw interface{}
}

type MapData map[string]interface{}
```

**RawData是原始数据；MapData是可以转换为map[string]interface{}的RawData；arr_data缓存数组；cache是缓存数据，有重合路径时，可以提高访问速度。**

##	核心函数

**Get**

```go
// Cache PathGet
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

// Cache Get
func (j *Jsnm) Get(path string) *Jsnm {
	if j==nil {
		return nil
	}
	// first step: get data from cache
	if nil == j.cache {
		j.cache = make(map[string]*Jsnm)
	} else {
		if cache_data, ok := j.cache[path]; ok {
			return cache_data
		}
	}
	// second step: get data from mapdata
	if j.map_data == nil {
		if map_data, ok := j.raw_data.(map[string]interface{}); ok {
			j.map_data = map_data
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
	j.cache[path] = will_cache_data
	return will_cache_data
}
```

-------------------------------

**Arr**

```go
// Cache Arr
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

// Cache ArrLocs
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

// Cache ArrLoc i
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
	j.arr_data = arr_cache
	j.arr_data[i] = NewJsnm(arr[i])
	if i >= len(arr) {
		return nil
	}
	return j.arr_data[i]
}

// Cache ArrPath path
func (j *Jsnm) ArrPath(path ...interface{}) *Jsnm {
	if len(path) <= 0 {
		return j
	}
	switch reflect.TypeOf(path[0]).Kind() {
	case reflect.String:
		return j.Get(path[0].(string)).ArrPath(path[1:]...)
	case reflect.Int:
		return j.ArrLoc(path[0].(int)).ArrPath(path[1:]...)
	}
	return nil
}
```

**具体的类型转换，可添加函数实现。**

_Example_

```go
age := jm.Get("Friends").Get("Age").MustInt64()
fmt.Println(age)
```


##	Benchmark

`go test -test.bench=".*"`

![Test][3]

 [1]: https://github.com/toukii/jsnm "jsnm"
 [2]: http://7xku3c.com1.z0.glb.clouddn.com/jsnm-benchmark.png "jsnm-bench"
 [3]: http://7xku3c.com1.z0.glb.clouddn.com/benchmark-jsnm.png "jsnm-bench"