package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		for v := range c {
			fmt.Println(v)
		}
	}()
	c <- 10
	c <- 100
	c <- 10
	c <- 0
	c <- -1234
	close(c)
	time.Sleep(time.Second)
	v, ok := <-c
	fmt.Println(v, ok)
}
