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
	return j.map_data
}
