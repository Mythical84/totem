package main

import (
	"fmt"
	"main/interpreter"
	"main/lexer"
	"main/parser"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Specify an input file")
		return
	}
	filename := os.Args[1]

	_, e := os.Stat(filename)
	if os.IsNotExist(e) {
		fmt.Println("File does not exist")
		return
	}

	data, _ := os.ReadFile(filename)
	content := string(data)
	tokens, err := lexer.Lexer(content, filename)
	if err != nil {
		for _, e := range err {
			os.Stderr.WriteString(e.Error() + "\n")
		}
		return
	}
	fmt.Printf("token stream: %+v\n", tokens)

	ast, err := parser.Parse(tokens, content, filename)
	if err != nil {
		for _, e := range err {
			os.Stderr.WriteString(e.Error() + "\n")
		}
		return
	}
	fmt.Printf("ast: %#v\n", ast)

	fmt.Println("\n---Program Output---")

	interpreter.Interpret(ast, filename);
}
