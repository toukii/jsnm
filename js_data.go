package jsnm

type Jsnm struct {
	raw   RawData
	data  MapData
	cache map[string]*Jsnm
}

func NewJsnm(m MapData) *Jsnm {
	return &Jsnm{
		raw:   RawData{raw: m},
		data:  m,
		cache: make(map[string]*Jsnm),
	}
}

func NewRawJsnm(raw interface{}) *Jsnm {
	return &Jsnm{
		raw:   NewRawData(raw),
		data:  nil,
		cache: nil,
	}
}

func (j *Jsnm) RawData() *RawData {
	if nil == j {
		return nil
	}
	return &j.raw
}

func (j *Jsnm) MapData() MapData {
	return j.data
}
