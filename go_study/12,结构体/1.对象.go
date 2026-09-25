package main

import "fmt"

//
//type Student struct {
//	Name string
//}
//
//func (s Student) Study() {
//	fmt.Printf("%s 正在学习", s.Name)
//}
//func main() {
//	s1 := Student{Name: "xiang"}
//	s1.Study()
//}

// 用 type 定义结构体 名字是Student 两个字段
type Student struct {
	Name string
	Age  int
}

// 值接收者方法 SetAge
func (s Student) SetAge(age int) {
	s.Age = age
}

// 指针接收者方法 SetAge1
func (s *Student) SetAge1(age int) {
	s.Age = age
}

// 5. main 函数
func main() {
	s := Student{
		Name: "XAIANG",
		Age:  24,
	}
	s.SetAge(18)
	fmt.Println(s.Age) // 输出 24，因为 SetAge 是值接收者
	s.SetAge1(18)
	fmt.Println(s.Age) // 输出 18，因为 SetAge1 是指针接收者
}
