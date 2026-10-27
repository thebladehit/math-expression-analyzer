package main

import (
	"bufio"
	"fmt"
	"os"
	an "pzcs-lab1/an"
	as "pzcs-lab1/associate"
	cm "pzcs-lab1/commute"
	pr "pzcs-lab1/parallel"
	pf "pzcs-lab1/postfix"
	sm "pzcs-lab1/simplifier"
	tr "pzcs-lab1/tree"
	"strings"
)

const formsLimit = 5000

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
	commuter := cm.NewCommuter()
	associator := as.NewAssociator()

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

	PrintForms("Форми за законом комутативності", commuter.Forms(tree, formsLimit), formsLimit)
	PrintForms("Форми за законом асоціативності", associator.Forms(tree, formsLimit), formsLimit)

	var both []*tr.Node
	for _, f := range commuter.Forms(tree, formsLimit) {
		both = append(both, associator.Forms(f, formsLimit)...)
	}
	both = tr.Unique(both)
	PrintForms("Форми за обома законами", both[:min(len(both), formsLimit)], formsLimit)

	tree = parallelizer.Parallelize(tree)
	PrintTree("Дерево паралельної форми", tree)
}
