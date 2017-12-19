package jsnm

import (
	"fmt"
)

type MapData map[string]interface{}

func (m *MapData) String() string {
	return fmt.Sprintf("{map_data:%#v}", m)
}

func (m MapData) Keys() []string {
	md := (map[string]interface{})(m)
	keys := make([]string, 0, len(md))
	for k, _ := range md {
		keys = append(keys, k)
	}
	return keys
}
