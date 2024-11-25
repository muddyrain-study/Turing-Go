package main

import (
	"fmt"
	"sort"
)

type People struct {
	Name string
	Age  int
}

type testSlice []People

func (l testSlice) Len() int           { return len(l) }
func (l testSlice) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l testSlice) Less(i, j int) bool { return l[i].Age < l[j].Age }

func main() {
	ls := testSlice{
		{Name: "n1", Age: 12},
		{Name: "n2", Age: 11},
		{Name: "n3", Age: 10},
	}

	fmt.Println(ls)
	sort.Sort(ls)
	fmt.Println(ls)
}
