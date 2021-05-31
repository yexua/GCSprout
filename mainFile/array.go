package main

import "fmt"

func main() {
	a := [...]int{1, 2, 3, 4}
	array(&a)
	fmt.Println(a[2])


	s := make([]byte, 5)
	fmt.Println(len(s))
	fmt.Println(cap(s))

	s = s[2:4]
	fmt.Println(len(s))
	fmt.Println(cap(s))

	s1 := []byte{'p', 'o', 'e', 'm'}

	s2 := s1[2:]
	fmt.Println(s2)
	s2[1] = 't'

	fmt.Printf("%c\n", s1[1])
	fmt.Printf("%c\n", s2[1])

}

func array(array *[4]int) {
	array[2] = 6
}
