package jsnm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/shaalx/goutils"
	gsj "github.com/shaalx/membership/pkg3/go-simplejson"
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

func TestPathGet(t *testing.T) {
	path_get := jm.PathGet("Friends", "One", "Name").RawData().String()
	assert(t, path_get, "One")
	path_get = jm.PathGet("Friends", "One", "Name").RawData().String()
	assert(t, path_get, "One")
}

func TestNCGet(t *testing.T) {
	cur := jm.NCGet("Friends")

	one_name := cur.NCGet("One", "Name")
	assert(t, one_name.RawData().String(), "One")

	one_name_X := jm.NCGet("Friends", "One", "Name", "X")
	if one_name_X != nil {
		t.Error(one_name_X, "should be nil.")
	}

	xx := one_name_X.NCGet("XX")
	if xx != nil {
		t.Error(xx, "should be nil.")
	}

	fon := jm.NCGet("Friends").NCGet("One").NCGet("Name")
	assert(t, fon.RawData().String(), "One")

	i64, _ := jm.NCGet("Age").RawData().Int64()
	assert(t, i64, int64(2))

	i64 = jm.NCGet("Age").RawData().MustInt64()
	assert(t, i64, int64(2))

}

func TestArr(t *testing.T) {
	arr := jm.Get("Loc").Arr()
	name := arr[0].RawData().String()
	assert(t, name, "Two")
	assert(t, arr[1].RawData().String(), "TwoTwo")

	arr1 := jm.Get("Loc").ArrLoc(1).RawData().String()
	assert(t, arr1, "TwoTwo")
}

func TestArr_NCGet(t *testing.T) {
	arr := jm.NCGet("Loc").Arr()
	name := arr[0].RawData().String()
	assert(t, name, "Two")
	assert(t, arr[1].RawData().String(), "TwoTwo")

	arr1 := jm.NCGet("Loc").ArrLoc(1).RawData().String()
	assert(t, arr1, "TwoTwo")
}

func TestArrJson(t *testing.T) {
	us := []*User{NewU("foo", 10), NewU("bar", 12)}
	bs, err := json.MarshalIndent(us, "\t", "\t")
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(string(bs))
	jmb := BytesFmt(bs)
	name := jmb.ArrLoc(0).Get("Name").RawData().String()
	assert(t, name, "foo")

	ioutil.WriteFile("foo.json", bs, 0666)
	jmf := FileNameFmt("foo.json")
	namef := jmf.ArrLoc(0).Get("Name").RawData().String()
	assert(t, namef, "foo")
}

func TestGsj(t *testing.T) {
	js, _ := gsj.NewJson(goutils.ReadFile("test.json"))
	name := js.GetPath("Friends", "One", "Name").MustString()
	assert(t, name, "One")
}

type U []*User
type S struct {
	U
}

func BenchmarkArr(b *testing.B) {
	b.StopTimer()
	// func TestBArr(b *testing.T) {
	U := []*User{NewU("O", 1), NewU("T", 2)}
	us := []S{S{U: U}, S{U: U}}
	bs, _ := json.MarshalIndent(us, "\t", "\t")
	// fmt.Println(string(bs))
	jmb := BytesFmt(bs)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = jmb.Arr()[0].Get("U").Arr()[0].Get("Loc").Arr()[0]
	}
	// ars := jmb.Arr()[0].Get("U").Arr()[0].Get("Loc").Arr()[0].RawData().String()
	// fmt.Println(ars)
}

func BenchmarkPathGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.PathGet("Friends", "One", "Name")
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.Get("Friends", "One", "Name")
	}
}

func BenchmarkShort_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.Get("Friends")
	}
}

func BenchmarkNCGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.NCGet("Friends", "One", "Name")
	}
}

func BenchmarkShort_NCGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.NCGet("Friends")
	}
}

func BenchmarkGsj(b *testing.B) {
	b.StopTimer()
	js, _ := gsj.NewJson(goutils.ReadFile("test.json"))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = js.GetPath("Friends", "One", "Name")
	}
}

func BenchmarkShort_Gsj(b *testing.B) {
	b.StopTimer()
	js, _ := gsj.NewJson(goutils.ReadFile("test.json"))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = js.GetPath("Friends")
	}
}
