package main

//go build -gcflags="-m" main.go

import "fmt"

type Data struct {
	Value int
}

//потенц. аллокация в куче (Pointer Semantics)
func CreatePointer() *Data {
	x := Data{Value: 42}
	return &x
}

//потенц. аллокация на стеке (Value Semantics)
func CreateValue() Data {
	y := Data{Value: 100}
	return y
}

func main() {
	p := CreatePointer()
	v := CreateValue()

	fmt.Println(p.Value, v.Value)
}