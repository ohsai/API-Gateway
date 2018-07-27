package main

import "fmt"

func new_thread(sigterm chan int) {
	for {
		switch sigval := <-sigterm; sigval {
		case 1:
			fmt.Println("SIGCONT")
		case 0:
			fmt.Println("SIGTERM")
			return
		}
	}
	close(sigterm)
}

func main() {
	signal := make(chan int)
	go new_thread(signal)
	fmt.Println("SIGCONT IN")
	signal <- 1
	fmt.Println("SIGCONT IN")
	signal <- 1
	fmt.Println("SIGTERM IN")
	signal <- 0
	fmt.Println("End of Program")
}
