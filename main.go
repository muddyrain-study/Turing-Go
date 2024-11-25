package main

import (
	"encoding/json"
	"os"
)

type Person struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Email  string   `json:"email"`
	Parent []string `json:"parent"`
}

func main() {
	p := Person{
		Name:   "zhangsan",
		Age:    20,
		Email:  "zhangsan@mail.com",
		Parent: []string{"Daddy", "Mom"},
	}
	f, _ := os.OpenFile("test.json", os.O_WRONLY, 077)
	defer f.Close()

	d := json.NewEncoder(f)
	d.Encode(p)
}
