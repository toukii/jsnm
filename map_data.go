package jsnm

import (
	"fmt"
)

type MapData map[string]interface{}

func (m *MapData) String() string {
	return fmt.Sprintf("{map_data:%#v}", m)
}
