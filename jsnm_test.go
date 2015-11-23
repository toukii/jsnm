package jsnm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	_ = ioutil.WriteFile(u.Name+".json", bs, 0666)
}

func NewU(n string, a byte) *User {
	u := &User{Name: n, Age: a}
	u.Friends = make(map[string]*User)
	u.Friends[n] = nil
	u.Loc = []string{n, n + n}
	return u
}

func MoreU() {
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
	jm = FileNameFmt("Two.json")
}

func TestGet(t *testing.T) {
	cur := jm.Get("Friends")
	fmt.Println(cur.RawData())

	one := cur.Get("One")
	fmt.Println(one.RawData())

	one_name := cur.Get("One", "Name")
	fmt.Println(one_name.RawData())

	one_name_X := jm.Get("Friends", "One", "Name", "X")
	fmt.Println(one_name_X)

	xx := one_name_X.Get("XX")
	fmt.Println(xx)

	fon := jm.Get("Friends").Get("One").Get("Name")
	fmt.Println(fon.RawData())

	i64, err := jm.Get("Age").RawData().Float64()
	fmt.Println(i64, err)
}

func TestArr(t *testing.T) {
	arr := jm.Get("Loc").Arr()
	fmt.Printf("%#v\n", arr)
	fmt.Println(arr[0].RawData())
	name := arr[0].Get("Name")
	fmt.Println(name)
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
