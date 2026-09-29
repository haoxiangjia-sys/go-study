# Go 学习笔记

个人 Go 语言学习过程中的练习代码与笔记，涵盖 Go 基础语法，以及 COMS20008 课程的 lab 练习。

## 目录结构

```
.
├── go_study/          # Go 基础语法练习（按学习顺序编号）
│   ├── 1.变量定义.go
│   ├── 2.输出.go
│   ├── 4.基本数据类型.go
│   ├── 5.数组.go
│   ├── 6.切片.go
│   ├── 8.map.go
│   ├── 9.if.go
│   ├── 10.switch.go
│   ├── 11.for.go
│   ├── 12,结构体/     # 结构体
│   ├── 14.接口.go
│   ├── 15.协程.go
│   ├── 16.channel.go
│   ├── 17.select.go
│   ├── 19.线程不安全.go
│   ├── 20.异常处理.go
│   ├── 21.泛式.go
│   ├── 函数/          # 简单函数 / 闭包 / 高阶函数
│   └── version/       # 版本常量示例
└── lab/               # 课程 lab 练习
    ├── week1/         # Intro to Go Lab 1 & 2
    └── week2/         # Distributed Lab & Lab 1
```

## 如何运行

`go_study/` 下的每个 `.go` 文件大多是独立的 `package main`，用 `go run` 单独运行：

```bash
go run go_study/1.变量定义.go
```

> ⚠️ 注意：这些练习文件各自都带 `func main()`，属于相互独立的程序，所以不能在一个目录下用 `go run .` 一起编译，否则会报 `multiple main functions`。

## 环境

- Go 1.x
- GoLand / VS Code
