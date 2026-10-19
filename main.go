package main

import (
	"bufio"
	"fmt"
	"os"
	an "pzcs-lab1/an"
	pr "pzcs-lab1/parallel"
	pf "pzcs-lab1/postfix"
	sm "pzcs-lab1/simplifier"
	tr "pzcs-lab1/tree"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	lexer := an.NewLexer()
	startState := an.StState
	parser := an.NewParser(startState)
	analyzer := an.NewAnalyzer(lexer, parser)
	postfixer := pf.NewPostfix()
	treeBuilder := tr.NewTreeBuilder()
	simplifier := sm.NewSimplifier()
	parallelizer := pr.NewParallelizer()

	fmt.Print("Введіть вираз: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка читання:", err)
		return
	}

	input = strings.Trim(input, "\n")

	tokens, errs := analyzer.Start(input)
	PrintResult(input, tokens, errs)
	if errs != nil {
		return
	}

	tokens = postfixer.Convert(tokens)
	fmt.Println("\nПостфіксна форма:", TokensToString(tokens))

	tree := treeBuilder.Build(tokens)
	PrintTree("Дерево виразу до оптимізацій", tree)

	tree, optimizations := simplifier.Simplify(tree)
	PrintOptimizations(optimizations)
	PrintTree("Дерево виразу після оптимізацій", tree)

	tree = parallelizer.Parallelize(tree)
	PrintTree("Дерево паралельної форми", tree)
}
