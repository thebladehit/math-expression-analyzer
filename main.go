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
	analyzer := an.NewAnalyzer()

	fmt.Print("Введіть вираз: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка читання:", err)
		return
	}

	input = strings.Trim(input, "\n")

	analyzer.Start(input)
	//res, errS := analyzer.Start(input)
	//if errS != nil {
	//	fmt.Println(errS)
	//}
	//fmt.Println(res)
}
