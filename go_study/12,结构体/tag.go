package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name     string `json:"name"` //tag
	Age      int
	Password string
}

func main() {
	user := User{Name: "xiang", Age: 20, Password: "123456"}
	byteData, _ := json.Marshal(user) //JSON 是为了跨语言、跨平台地传输和存储数据。
	fmt.Println(string(byteData))
}
