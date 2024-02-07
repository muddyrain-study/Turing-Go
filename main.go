package main

import (
	"log"
	"strings"
)

func main() {
	s := "Hello,码神之路Java教程"
	source := "Java"
	target := "Go"

	index := strings.Index(s, "Java")
	sourceLen := len(source)

	start := s[:index]
	log.Default().Println(start)
	end := s[index+sourceLen:]
	log.Default().Println(end)

	s = start + target + end
	log.Default().Println(s)
}
