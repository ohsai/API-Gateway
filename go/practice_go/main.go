package main

import "fmt"

func sum(a []int, c chan int) {
	for _, v := range a {
		fmt.Println(v)
		c <- v
	}
	close(c)
}

func main() {
	a := []int{1, 2, 3, 4, 5}
	c := make(chan int)
	go sum(a, c)
	for i := range c {
		fmt.Println(i)
	}
}
