package player

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
