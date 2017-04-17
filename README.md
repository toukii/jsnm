#	[Jsnm][1]

---------------------


__json mapping for map[string]interface{}__

提供缓存的json解析器，初衷是提高重复解析json速度。


##	数据结构设计

```go
type Jsnm struct {
	raw   RawData
	data  MapData
	cache map[string]*Jsnm
}

type RawData struct {
	raw interface{}
}

type MapData map[string]interface{}
```

**RawData是原始数据；MapData是可以转换为map[string]interface{}的RawData；cache是缓存数据，有重合路径时，可以提高访问速度。**

##	核心函数

**Get**

```go
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

// Get data first from cache data
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
```
_Example_

```go
jm.Get("Friends").Get("Age")
// NCGet should be after the Get
fon := jm.Get("Friends").NCGet("One").NCGet("Name")
```

-------------------------------

**Arr**

```go
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

// Get the i-th index from array
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
```
_Example_

```go
jm.Get("Loc").Arr[0].Get("Name")
arr1 := jm.NCGet("Loc").ArrLoc(1).RawData().String()
```

**具体的类型转换，可在RawData中添加函数实现。**

_Example_

```go
age := jm.Get("Friends").Get("Age").MustInt64()
fmt.Println(age)
```


##	Benchmark

`go test -test.bench=".*"`

![Test][2]

 [1]: https://github.com/shaalx/jsnm "jsnm"
 [2]: http://7xku3c.com1.z0.glb.clouddn.com/jsnm-benchmark.png "jsnm-bench"