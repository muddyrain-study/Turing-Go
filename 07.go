package main

func playerGen(name string) func() (string, int) {
	// 血量一直为150
	hp := 150
	return func() (string, int) {
		return name, hp
	}
}
func playerGen1() func(string) (string, int) {
	// 血量一直为150
	hp := 150
	return func(name string) (string, int) {
		return name, hp
	}
}
func main() {
	generator := playerGen1()
	name1, hp1 := generator("玩家1")
	name2, hp2 := generator("玩家2")
	name3, hp3 := generator("玩家3")
	println(name1, hp1)
	println(name2, hp2)
	println(name3, hp3)
}
