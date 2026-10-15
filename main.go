package main

import (
	"fmt"
	"os"
)

func main() {
	input := "sample.txt"
	if len(os.Args) > 1 {
		input = os.Args[1]
	}
	if len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . [entree.txt] [sortie.txt]")
		os.Exit(1)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur de lecture :", err)
		os.Exit(1)
	}

	result := Process(string(data))

	if len(os.Args) == 3 {
		if err := os.WriteFile(os.Args[2], []byte(result), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "erreur d'écriture :", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println("---Texte original---")
	fmt.Print(string(data))
	fmt.Println("\n--- Texte corrigé ---")
	fmt.Println(result)
}
