package main

import (
	"bufio"
	"fmt"
	"os"
	an "pzcs-lab1/an"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	lexer := an.NewLexer()
	startState := an.StState
	parser := an.NewParser(startState)
	analyzer := an.NewAnalyzer(lexer, parser)

	fmt.Print("Введіть вираз: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка читання:", err)
		return
	}

	input = strings.Trim(input, "\n")

	tokens, errs := analyzer.Start(input)
	fmt.Println(tokens)
	if errs != nil {
		fmt.Println(errs)
	}
}
