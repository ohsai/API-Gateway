package main

import (
	"fmt"
	"time"
)

func new_thread(sigchld chan int, index int) {
	time.Sleep(time.Duration(index) * time.Second)
	fmt.Println(index, "-th thread terminated")
	sigchld <- 1
	return
}

func main() {
	sigchld := make(chan int, 2)
	number_of_cores := 8
	for index := 0; index < number_of_cores; index++ {
		go new_thread(sigchld, index)
	}
	for index := 0; index < number_of_cores; index++ {
		x := <-sigchld
		fmt.Println("received sigchld status : ", x)
	}

	/*
		sigchldlist := make ([]chan int , number_of_cores )
		for index, value := range sigchldlist {
			fmt.Println(index, "-th thread created")
			go new_thread(value)
		}
	*/
	/*
		fmt.Println("SIGCONT IN")
		signal <- 1
		fmt.Println("SIGCONT IN")
		signal <- 1
		fmt.Println("SIGTERM IN")
		signal <- 0
	*/
	fmt.Println("End of Program")
}
