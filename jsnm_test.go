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

func TestOne(t *testing.T) {
	// MoreU()

	jm := FileNameFmt("Two.json")
	fmt.Println(jm)
	v, err := jm.Get("Friend")
	fmt.Println(v, err)

	v, err = jm.Get("Friends", "One", "Loc")
	fmt.Println(v, err)

	v, err = jm.Get("Friends", "One", "Name")
	fmt.Println(v, err)

	v, err = jm.Get("Friends", "One", "Name")
	fmt.Println(v, err)

	fmt.Println(jm)
}
