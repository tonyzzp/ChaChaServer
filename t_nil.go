package main

import "fmt"

type s struct {
	v int
}

func build() *s {
	return nil
}

func adder(i *int) {
	*i++
}

func main() {
	i := 10
	fmt.Println(i)
	adder(&i)
	fmt.Println(i)
}
