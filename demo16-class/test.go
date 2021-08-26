package main

import "fmt"

type I interface {
	Say()
	Run()
}

type XXX struct {
	Name string
}

func (x XXX) Say() {
	fmt.Printf("XXX name is: %s\n", x.Name)
}

func (x XXX) Run() {
	fmt.Printf("XXX run away\n")
}

type YYY struct {
	XXX
	Age int
}

func (y YYY) Say() {
	fmt.Printf("YYY name is: %s, age is: %d\n", y.Name, y.Age)
}

func Play(p I) {
	p.Say()
	p.Run()

	_, ok := p.(*XXX)
	if !ok {
		fmt.Printf("not xxx\n")
	} else {
		fmt.Printf("is xxx\n")
	}
}

func main() {
	x := &XXX{"x-aka-x"}
	y := &YYY{XXX: *x, Age: 23}

	fmt.Printf("Play x\n")
	Play(x)
	fmt.Printf("Play y\n")
	Play(y)
}
