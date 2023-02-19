package model

import (
	"fmt"
	"testing"
)

type A struct {
	A string
}

type B struct {
	B string
}

func TestStruct(t *testing.T) {
	var x interface{} = &A{A: "123"}
	y, ok := x.(*A)
	fmt.Println(y, ok)
	z, ok := x.(*B)
	fmt.Println(z, ok)
}

func TestMedia_Copy(t *testing.T) {
	m := &Media{Title: "asdf", User: &User{Name: "123"}}
	m2 := m.Copy()
	fmt.Println(m, m2)
	m2.User.(*User).Name = "456"
	fmt.Println(m.User.(*User).Name, m2)
}
