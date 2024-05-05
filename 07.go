package main

func main() {
	var whatever = [5]int{1, 2, 3, 4, 5}
	for i := 0; i < len(whatever); i++ {
		i := i
		defer func() {
			println(i)
		}()
	}
}
