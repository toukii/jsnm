package jsnm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

type User struct {
	Name    string
	Age     byte
	Friends map[string]*User
	Loc     []string
}

func (u *User) ToJ() {
	bs, _ := json.MarshalIndent(u, "\t", "\t")
	_ = ioutil.WriteFile("test.json", bs, 0666)
}

func NewU(n string, a byte) *User {
	u := &User{Name: n, Age: a}
	u.Friends = make(map[string]*User)
	u.Friends[n] = nil
	u.Loc = []string{n, n + n}
	return u
}

func Mock() {
	u1 := NewU("One", 1)
	u2 := NewU("Two", 2)
	u2.Friends["One"] = u1
	u2.ToJ()
}

var (
	jm *Jsnm
)

func init() {
	fmt.Println("test...")
	Mock()
	jm = FileNameFmt("test.json")
}

func assert(t *testing.T, get, want interface{}) bool {
	if get == nil && want == nil {
		return true
	}
	if !reflect.DeepEqual(want, get) {
		t.Errorf("want:%v, get: %v\n", want, get)
		return false
	}
	return true
}

func TestGet(t *testing.T) {
	cur := jm.Get("Friends")

	one_name := cur.Get("One", "Name")
	assert(t, one_name.RawData().String(), "One")

	one_name_X := jm.Get("Friends", "One", "Name", "X")
	if one_name_X != nil {
		t.Error(one_name_X, "should be nil.")
	}

	xx := one_name_X.Get("XX")
	if xx != nil {
		t.Error(xx, "should be nil.")
	}

	fon := jm.Get("Friends").Get("One").Get("Name")
	assert(t, fon.RawData().String(), "One")

	i64, _ := jm.Get("Age").RawData().Int64()
	assert(t, i64, int64(2))

	i64 = jm.Get("Age").RawData().MustInt64()
	assert(t, i64, int64(2))

}

func TestArr(t *testing.T) {
	arr := jm.Get("Loc").Arr()
	name := arr[0].RawData().String()
	assert(t, name, "Two")
	assert(t, arr[1].RawData().String(), "TwoTwo")
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.Get("Friends", "One", "Name")
	}
}

func BenchmarkGetShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.Get("Friends")
	}
}
