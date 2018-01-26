package jsnm

import (
	"encoding/json"
	"fmt"
	//gsj "github.com/bitly/go-simplejson"
	"github.com/toukii/goutils"
	gsj "github.com/toukii/membership/pkg3/go-simplejson"
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
	mp := make(map[string]interface{})
	mp["objs"] = []*User{u1, u2}
	bs, _ := json.MarshalIndent(mp, "\t", "\t")
	_ = ioutil.WriteFile("test.json", bs, 0666)
	u2.Friends["One"] = u1
	u2.ToJ()
}

var (
	bs = []byte(`
[
	{
		"1":[
			{
				"Name": "foo",
				"Age": 10,
				"Friends": {
					"foo": {
						"11":[
							{
								"Name": "foo11",
								"Age": 1011,
								"Friends": {
									"foo11": 
										{
											"Name": "bar211",
											"Age": 1221,
											"Friends": {
												"bar21": null
											},
											"Loc": [
												"bar21",
												"barbar21"
											]
										}
								},
								"Loc": [
									"foo2",
									"foofoo2"
								]
							},
							{
								"Name": "bar2",
								"Age": 122,
								"Friends": {
									"bar2": null
								},
								"Loc": [
									"bar2",
									"barbar2"
								]
							}
						]
					}
				},
				"Loc": [
					"foo",
					"foofoo"
				]
			},
			{
				"Name": "bar",
				"Age": 12,
				"Friends": {
					"bar": null
				},
				"Loc": [
					"bar",
					"barbar"
				]
			}
		]
	}
]`)

	arr_bs = []byte(`[
    [
        [
            [
                "a"
            ]
        ]
    ]
]`)

	jm    *Jsnm
	jmGSJ *gsj.Json

	arrjm    *Jsnm
	arrjmGSJ *gsj.Json
)

func init() {
	fmt.Println("test...")
	Mock()
	jm = FileNameFmt("test.json")
	jmGSJ, _ = gsj.NewJson(goutils.ReadFile("test.json"))

	arrjm = BytesFmt(arr_bs)
	arrjmGSJ, _ = gsj.NewJson(arr_bs)
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

	one_name := cur.Get("One").Get("Name")
	assert(t, one_name.RawData().String(), "One")

	one_name_X := jm.Get("Friends").Get("One").Get("Name").Get("X")
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

func TestArrRange(t *testing.T) {
	loc := jm.Get("Loc")
	names := make([]string, 0, 10)
	loc.Range(func(i int, ji *Jsnm) {
		names = append(names, fmt.Sprintf("%d-%s", i, ji.RawData().String()))
	})

	assert(t, names, []string{"0-Two", "1-TwoTwo"})
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
	foo := jmf.ArrGet("0", "Name").RawData().String()
	assert(t, foo, "foo")
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

func TestLongArr(t *testing.T) {
	// // func TestBArr(b *testing.T) {
	// U := []*User{NewU("O", 1), NewU("T", 2)}
	// us := []S{S{U: U}, S{U: U}}
	// bs, _ := json.MarshalIndent(us, "\t", "\t")
	// fmt.Println(string(bs))
	jmb := BytesFmt(bs)
	// rs := jmb.Arr()[0].Get("1").PathGet("Friends","foo","11")//.Arr()[0].Get("Loc").Arr()[0]
	//rs := jmb.Arr()[0].Get("1").ArrLoc(0).Get("Friends").Get("foo").Get("11").Arr()[0].Get("Loc").Arr()[0]
	//rs := jmb.ArrGet("0",`"1"`,"0","Friends","foo",`"11"`,"0","Loc","0")
	rs := jmb.ArrPath(0, "1", 0, "Friends", "foo", "11", 0, "Loc", 0).String()
	//()[0].Get("1").ArrLoc(0).Get("Friends").Get("foo").Get("11").Arr()[0].Get("Loc").Arr()[0]
	assert(t, rs, "foo2")
	// fmt.Println(ars)

	jmb2, err := gsj.NewJson(bs)
	if goutils.CheckErr(err) {
		t.Error(err)
	}
	rs2 := jmb2.GetIndex(0).Get("1").GetIndex(0).GetPath("Friends", "foo", "11").GetIndex(0).Get("Loc").GetIndex(0).MustString()
	assert(t, rs2, "foo2")
}

func Benchmark_001_short_Get_jsnm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.Get("Friends")
	}
}

func Benchmark_001_short_NCGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.NCGet("Friends")
	}
}

func Benchmark_001_short_Get_gsj(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jmGSJ.Get("Friends")
	}
}

func Benchmark_002_PathGet_jsnm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.PathGet("Friends", "One", "Name")
	}
}

func Benchmark_002_NCGet_path(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jm.NCGet("Friends", "One", "Name")
	}
}

func Benchmark_002_GetPath_gsj(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = jmGSJ.GetPath("Friends", "One", "Name")
	}
}

func Benchmark_003_Get_jsnm(b *testing.B) {
	//b.StopTimer()
	//fmt.Println("PathGet:",jm.PathGet("Friends", "One", "Name").RawData().String())
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		//_= jm.PathGet("Friends", "One", "Name").RawData().String()
		_ = jm.Get("Friends").Get("One").Get("Name")
	}
}

func Benchmark_003_NCGet(b *testing.B) {
	//b.StopTimer()
	//fmt.Println("NCGET-:",jm.NCGet("Friends").NCGet("One").NCGet("Name").RawData().String())
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = jm.NCGet("Friends").NCGet("One").NCGet("Name")
	}
}

func Benchmark_003_Get_gsj(b *testing.B) {
	//b.StopTimer()
	//fmt.Println("PathGet-gsj:",jmGSJ.GetPath("Friends", "One", "Name"))
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = jmGSJ.Get("Friends").Get("One").Get("Name")
	}
}

func Benchmark_004_Arr_jsnm(b *testing.B) {
	//b.StopTimer()
	//fmt.Println("Arr:",arrjm.Arr()[0].Arr()[0].Arr()[0].Arr()[0])
	//fmt.Println("Arr:",arrjm.ArrLoc(1).ArrLoc(0).ArrLoc(0).ArrLoc(0))
	//fmt.Println("arrlocs:",arrjm.ArrLocs(0,0,0,0))
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		//_= arrjm.Arr()[0].Arr()[0].Arr()[0].Arr()[0]
		_ = arrjm.ArrLoc(0).ArrLoc(0).ArrLoc(0).ArrLoc(0)
		//_=arrjm.ArrLocs(0,0,0,0)
		//_=arrjm.ArrLoc(0).Arr()[0].ArrLoc(0).Arr()[0]
	}
}

func Benchmark_004_Arr_gsj(b *testing.B) {
	//b.StopTimer()
	//fmt.Println("Arr-gsj:",arrjmGSJ.GetIndex(0).GetIndex(0).GetIndex(0).GetIndex(0))
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = arrjmGSJ.GetIndex(0).GetIndex(0).GetIndex(0).GetIndex(0)
	}
}

func Benchmark_005_jsnm(b *testing.B) {
	b.StopTimer()
	jmb := BytesFmt(bs)
	//rs:= jmb.Arr()[0].Get("1").ArrLoc(0).PathGet("Friends","foo","11").Arr()[0].Get("Loc").Arr()[0]
	//fmt.Println(rs)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = jmb.Arr()[0].Get("1").ArrLoc(0).PathGet("Friends", "foo", "11").Arr()[0].Get("Loc").Arr()[0]
		//_=jmb.ArrGet("0",`"1"`,"0","Friends","foo",`"11"`,"0","Loc","0")
		//_= jmb.ArrPath(0,"1",0,"Friends","foo","11",0,"Loc",0)
	}
}

func Benchmark_005_gsj(b *testing.B) {
	b.StopTimer()
	jmb2, err := gsj.NewJson(bs)
	if goutils.CheckErr(err) {
		b.Error(err)
	}
	//rs := jmb2.GetIndex(0).Get("1").GetIndex(0).GetPath("Friends","foo","11").GetIndex(0).Get("Loc").GetIndex(0)
	//fmt.Print(rs)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = jmb2.GetIndex(0).Get("1").GetIndex(0).GetPath("Friends", "foo", "11").GetIndex(0).Get("Loc").GetIndex(0)
	}
}

func Benchmark_006_new_jsnm(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = &Jsnm{raw_data: bs}
	}
}

func Benchmark_006_new_gsj(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = gsj.NewJson(bs)
	}
}
