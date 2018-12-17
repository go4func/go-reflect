package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type outExp struct {
	out []reflect.Value
	exp time.Time
}

func Cacher(f interface{}, expiration time.Duration) (interface{}, error) {
	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func {
		return nil, errors.New("Only accept function.")
	}
	if ft.NumIn() == 0 {
		return nil, errors.New("must have at least one param.")
	}
	if ft.NumOut() == 0 {
		return nil, errors.New("must have at least one return value.")
	}
	inType, err := BuildInStruct(ft)
	if err != nil {
		return nil, err
	}
	fmt.Println("inType look like", inType)

	m := map[interface{}]outExp{}
	fv := reflect.ValueOf(f)
	cacher := reflect.MakeFunc(ft, func(in []reflect.Value) []reflect.Value {
		iv := reflect.New(inType).Elem()
		for i, v := range in {
			iv.Field(i).Set(v)
		}

		ov, ok := m[iv.Interface()]

		if !ok || ov.exp.Before(time.Now()) {
			ov.out = fv.Call(in)
			ov.exp = time.Now().Add(expiration)
			m[iv.Interface()] = ov
		}

		return ov.out

	})

	return cacher.Interface(), nil
}

func BuildInStruct(ft reflect.Type) (reflect.Type, error) {
	sfs := []reflect.StructField{}
	for i := 0; i < ft.NumIn(); i++ {
		f := ft.In(i)
		if !f.Comparable() {
			return nil, fmt.Errorf("parameter %d of type %s and kind %s is not comparable.",
				i+1, f.Name(), f.Kind())
		}
		sf := reflect.StructField{
			Name: fmt.Sprintf("Fields%d", i+1),
			Type: f,
		}
		sfs = append(sfs, sf)
	}

	s := reflect.StructOf(sfs)
	return s, nil

}
