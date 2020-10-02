package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		for i := 0; i <= 5; i++ {
			a <- i
		}
	}()

	go func() {
		for i := 10; i > 5; i++ {
			b <- i
		}
	}()

	for i := 0; i < 4; i++ {
		select {
		case v1 := <-a:
			fmt.Printf("%d\n", v1)
		case v2 := <-b:
			fmt.Printf("%d\n", v2)
		}

		time.Sleep(3 * time.Second)
		fmt.Println("hello")
	}
}
