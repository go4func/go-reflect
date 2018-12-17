package main

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

type T struct {
	Name string `json:"name" validate:"required"`
	P    *int   `json:"p" validate:"required"`
}

// .Type
func Examiner(t reflect.Type) {
	fmt.Println("Kind", t.Kind())
	fmt.Println("Type", t.Name())
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println("\tName:", f.Name, "\tKind:", f.Type.Kind(), "\tType:", f.Type, "\tTag:", f.Tag)
		}
	}
	//
}

func (t T) Clone() T {
	t2 := T{}
	vT := reflect.ValueOf(t)
	vT2 := reflect.ValueOf(&t2).Elem()
	for i := 0; i < vT.NumField(); i++ {
		f := vT.Field(i)
		f2 := vT2.Field(i)
		switch f.Kind() {
		case reflect.String:
			f2.Set(f)
		case reflect.Ptr:
			// if f2.IsNil() {
			// }
			// f2.Elem().Set(f.Elem())
			f2.Set(reflect.New(f.Type().Elem()))
			f2 = f.Elem()
		}
	}
	return t2
}

func main() {
	// c := &check{}
	// log.Println(ImplementChecker(reflect.TypeOf(c)))
	// t := T{Name: "old name", P: new(int)}
	// *t.P = 11
	// t2 := t.Clone()
	// *t2.P = 20
	// log.Println(*t.P)
	// log.Println(*t2.P)
	// c, err := Cacher(Adder, time.Duration(2)*time.Second)
	// ch := c.(func(int, int) int)
	// if err != nil {
	// 	panic(err)
	// }
	// for i := 0; i < 5; i++ {
	// 	log.Println("attemp #", i+1)
	// 	start := time.Now()
	// 	ch(1, 2)
	// 	log.Println("enlapse time:", time.Now().Sub(start))
	// }

	// // after sleep cache should be expired
	// // enlapse time back to normal
	// time.Sleep(time.Duration(3) * time.Second)
	// log.Println("attemp #", 6)
	// start := time.Now()
	// ch(1, 2)
	// log.Println("enlapse time:", time.Now().Sub(start))
	start := time.Now()
	time.Sleep(1 * time.Second)
	fmt.Println("enlapse time:", time.Now().Sub(start))
}

type Checker interface {
	DoWork()
}

type check struct {
}

func (c *check) DoWork() {
	log.Println("doing important jobs")
}

func ImplementChecker(tt reflect.Type) bool {
	it := reflect.TypeOf((*Checker)(nil)).Elem()
	return tt.Implements(it)
}

func CreateStruct(v ...interface{}) interface{} {
	var sfs []reflect.StructField
	for i, f := range v {
		rt := reflect.TypeOf(f)
		sf := reflect.StructField{
			Name: fmt.Sprintf("Field%d", i+1),
			Type: rt,
		}
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	sp := reflect.New(st)
	return sp.Interface()

}

func Adder(a, b int) int {
	// time.Sleep(time.Duration(200) * time.Millisecond)
	log.Printf("Adder: %d + %d = %d\n", a, b, a+b)
	return a + b
}

func MakeTimer(f interface{}) interface{} {
	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func {
		panic("invalid type, expected funtion")
	}

	v := reflect.ValueOf(f)
	timer := reflect.MakeFunc(ft, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := v.Call(in)
		log.Printf("enlapse time: %v\n", time.Now().Sub(start))
		return out
	})
	return timer.Interface()
}

func CreatePrimative(t reflect.Type) reflect.Value {
	return reflect.Zero(t)
	// chType:= reflect.MakeChan(reflect.ChanOf(reflect.BothDir, t),10)
	// v := chType.Interface()
	// reflect.New()
}
