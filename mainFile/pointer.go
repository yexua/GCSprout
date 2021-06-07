package main

import (
	"fmt"
)
func main() {
	var i1 = 5
	fmt.Printf("An integer: %d, its location in memory: %p\n", i1, &i1)
	var intP *int
	intP = &i1
	fmt.Printf("The value at memory location %p is %d\n", intP, *intP)

	ptr := &i1
	fmt.Println(ptr)

	var num1 int = 7

	switch {
	case num1 < 0:
		fmt.Println("Number is negative")
	case num1 > 0 && num1 < 10:
		fmt.Println("Number is between 0 and 10")

	default:
		fmt.Println("Number is 10 or greater")
	}



	for i := 0; i<7 ; i++ {
		if i%2 == 0 { continue }
		fmt.Println("Odd:", i)
	}


}
