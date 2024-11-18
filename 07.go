package main

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func main() {
	//Try(func() {
	//	panic("test panic")
	//}, func(err interface{}) {
	//	fmt.Println(err)
	//})

	//b := []int{0, 0, 0}
	//a(b)
	//for i, i2 := range b {
	//	println(i, i2)
	//}

	//p1 := player()
	//p2 := player()
	//fmt.Println(p1())
	//fmt.Println(p1())
	//fmt.Println(p1())
	//
	//fmt.Println(p2())

	//start := time.Now()
	//fmt.Println(start)

	//defer func() {
	//	fmt.Println("defer1")
	//	if err := recover(); err != nil {
	//		println("1", err.(string))
	//	}
	//}()
	//
	//defer func() {
	//	fmt.Println("defer2")
	//	if err := recover(); err != nil {
	//		println("2", err.(string))
	//	}
	//}()
	//panic("test panic")
	try(func() {
		panic("test panic")
	}, func(err interface{}) {
		println(err.(string))
	})
}

func try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}
func a(s []int) {
	s[0] = 100
}

func player() func() int {
	hp := 30

	return func() int {
		hp++
		return hp
	}
}
