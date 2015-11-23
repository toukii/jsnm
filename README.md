#	jsnm

**json mapping for map[string]interface{}**

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
```
_Example_

```go
jm.Get("Friends").Get("Age")
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
```
_Example_

```go
jm.Get("Loc").Arr[0].Get("Name")
```

**具体的类型转换，可在RawData中添加函数实现。**

_Example_

```go
age := jm.Get("Friends").Get("Age").MustInt64()
fmt.Println(age)
```


 [1]: https://github.com/shaalx/jsnm "jsnm"