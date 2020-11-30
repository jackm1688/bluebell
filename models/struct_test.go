package models

import (
	"fmt"
	"testing"
	"unsafe"
)

//内存对齐
type s1 struct {
	a int8
	b string
	c int8
}

type s2 struct {
	a int8
	c int8
	b string
}

func TestStruct(t *testing.T) {
	v1 := s1{
		a: 1,
		b: "hello",
		c: 2,
	}

	v2 := s2{
		a: 1,
		c: 2,
		b: "hello",
	}

	fmt.Println(unsafe.Sizeof(v1), unsafe.Sizeof(v2))

	fmt.Println()
	var s = []int{1, 2, 3, 4, 5}
	//var b = []int{4, 5, 6}
	//copy(b, s) //将s复制到b

	/*fmt.Println(s[1 : len(s)-1])
	fmt.Println(s[2:])
	copy(s[1:len(s)-1], s[2:])*/
	s = append(s[:1], s[2:]...)

	fmt.Println(s)

}
