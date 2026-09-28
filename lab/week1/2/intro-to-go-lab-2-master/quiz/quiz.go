package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type question struct {
	q, a string
}

type score int

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func questions() []question {
	f, err := os.Open("quiz-questions.csv")
	check(err)
	reader := csv.NewReader(f)
	table, err := reader.ReadAll()
	check(err)
	var questions []question
	for _, row := range table {
		questions = append(questions, question{q: row[0], a: row[1]})
	}
	return questions
}

// 4a: ask 现在接收一个 result channel，用于把分数传回 main
func ask(s score, question question, result chan score) {
	fmt.Println(question.q)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter answer: ")
	scanner.Scan() // 这里会阻塞，等待用户输入
	text := scanner.Text()
	if strings.Compare(text, question.a) == 0 {
		fmt.Println("Correct!")
		s++
	} else {
		fmt.Println("Incorrect :-(")
	}
	result <- s // 将更新后的分数发送回 main
}

func main() {
	s := score(0)
	qs := questions()
	result := make(chan score)
	timer := time.After(5 * time.Second) // 创建一个 5 秒的定时器

	for _, q := range qs {
		go ask(s, q, result) // 启动 goroutine 等待用户输入

		select {
		case newScore := <-result:
			s = newScore // 用户回答了，更新分数，准备下一题
		case <-timer:
			// 5秒时间到，立即终止
			fmt.Println("\nTime's up!")
			fmt.Println("Final score", s)
			return // 直接退出 main，不等待当前 ask goroutine
		}
	}

	// 如果所有题目都在 5 秒内回答完了，也会走到这里
	fmt.Println("Final score", s)
}
