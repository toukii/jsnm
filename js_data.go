package jsnm

type Jsnm struct {
	raw_data interface{}
	map_data MapData
	arr_data []*Jsnm
	cache    map[string]*Jsnm
}

func NewJsnm(raw interface{}) *Jsnm {
	return &Jsnm{raw_data: raw}
}

func (j *Jsnm) RawData() *RawData {
	if nil == j {
		return nil
	}
	return &RawData{raw: j.raw_data}
}

func (j *Jsnm) MapData() MapData {
	if j.map_data == nil {
		j.map_data = make(MapData)
		if map_data, ok := j.raw_data.(map[string]interface{}); ok {
			j.map_data = map_data
		} else {
			return nil
		}
	}

	return j.map_data
}

func (j *Jsnm) MustFloat64() float64 {
	return j.RawData().MustFloat64()
}

func (j *Jsnm) MustInt64() int64 {
	return j.RawData().MustInt64()
}
