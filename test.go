package main

import "fmt"

type Controller struct {
	paths []int
}

func main() {
	a := []int{1, 2, 3}
	b := append(a, 4)
	controller := Controller{
		paths: b,
	}
	b[0] = 100
	// c := append(a, 5)
	fmt.Println(controller)

}
