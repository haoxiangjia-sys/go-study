package main

import "fmt"

// Animal 定义一个animal的接口，它有唱，跳，rap的方法
type Animal interface {
	sing()
	jump()
	rap()
}

// Chicken 需要全部实现这些接口
type Chicken struct {
	Name string
}

func (c Chicken) sing() {
	fmt.Println("chicken 唱")
}

func (c Chicken) jump() {
	fmt.Println("chicken 跳")
}
func (c Chicken) rap() {
	fmt.Println("chicken rap")
}

// Cat 需要全部实现这些接口
type Cat struct {
	Name string
}

func (c Cat) sing() {
	fmt.Println("cat 唱")
}

func (c Cat) jump() {
	fmt.Println("cat 跳")
}
func (c Cat) rap() {
	fmt.Println("cat rap")
}

func sing(obj Animal) {
	obj.sing()
}

// 全部实现完之后，chicken就不再是一只普通的鸡了

func main() {
	chicken := Chicken{"ik"}
	cat := Cat{"阿狸"}
	sing(chicken)
	sing(cat)
}
