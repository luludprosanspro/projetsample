package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	fichierEntree := "sample.txt"
	if len(os.Args) > 1 {
		fichierEntree = os.Args[1]
	}

	contenu, err := os.ReadFile(fichierEntree)
	if err != nil {
		fmt.Println("Impossible de lire le fichier :", err)
		return
	}

	texteCorrige := Process(string(contenu))

	if len(os.Args) == 3 {
		err = os.WriteFile(os.Args[2], []byte(texteCorrige), 0o644)
		if err != nil {
			fmt.Println("Impossible d'écrire le fichier :", err)
		}
		return
	}
