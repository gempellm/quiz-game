package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type stage struct {
	question string
	answer   string
}

func main() {
	stages := fillStages()
	startQuiz(stages)
}

func fillStages() []stage {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	length := 10
	maxVal := 50

	stages := make([]stage, length, length)

	for i := 0; i < length; i++ {
		x := rand.Intn(maxVal)
		y := rand.Intn(maxVal)
		z := x + y

		stages[i] = stage{question: fmt.Sprintf("%v + %v", x, y), answer: fmt.Sprint(z)}
	}

	return stages
}

func startQuiz(stages []stage) {
	fmt.Println("~~ Math Quiz ~~")

	reader := bufio.NewReader(os.Stdin)
	total := len(stages)
	correct := 0

	ch := make(chan string)

	// When the main thread returns (i.e., you return from the main function), it terminates the entire process.
	// Это значит, что данный поток не уйдет в бесконечное ожидание ввода строки, а завершится вместе с main()
	// Однако при этом данная конструкция позволяет не уходить в бесконечное ожидание ввода строки для старой переменной
	// при срабатывании таймера, а оно считывает входные данные для актуальной переменной stages
	go func(ch chan string) {
		for {
			readAnswer(reader, ch)
		}
	}(ch)

	for _, v := range stages {
		timer := make(chan bool)
		defer close(timer)
		fmt.Print(v.question, " = ")

		go startTimer(timer)

		select {
		case a := <-ch:
			a = strings.TrimSpace(a)
			if a == v.answer {
				correct++
			}
		case <-timer:
			fmt.Println()
			continue
		}
	}

	fmt.Println("Correct answers:", correct, "/", total)
}

func readAnswer(reader *bufio.Reader, ch chan string) {
	answer, _ := reader.ReadString('\n')
	ch <- answer
}

func startTimer(ch chan bool) {
	time.Sleep(3 * time.Second)
	ch <- true
}
