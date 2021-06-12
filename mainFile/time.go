package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "我爱你"
	fmt.Println(utf8.RuneCountInString(s))
}
