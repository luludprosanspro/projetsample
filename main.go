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

	fmt.Println("--- Texte original ---")
	fmt.Print(string(contenu))
	fmt.Println("\n--- Texte corrigé ---")
	fmt.Println(texteCorrige)
}

func Process(texte string) string {
	lignes := strings.Split(texte, "\n")
	for i, ligne := range lignes {
		lignes[i] = corrigerLigne(ligne)
	}
	return strings.Join(lignes, "\n")
}

func corrigerLigne(ligne string) string {
	mots := decouper(ligne)
	mots = appliquerBalises(mots)
	corrigerAAn(mots)
	return recomposer(mots)
}


func decouper(ligne string) []string {
	lettres := []rune(ligne)
	var mots []string
	motActuel := ""

	finirMot := func() {
		if motActuel != "" {
			mots = append(mots, motActuel)
			motActuel = ""
		}
	}

	for i := 0; i < len(lettres); i++ {
		c := lettres[i]

		if unicode.IsSpace(c) {
			finirMot()

		} else if c == '(' {

			fin := trouverParenthese(lettres, i)
			if fin != -1 && estBalise(string(lettres[i:fin+1])) {
				finirMot()
				mots = append(mots, string(lettres[i:fin+1]))
				i = fin 
			} else {
				motActuel += string(c)
			}

		} else if strings.ContainsRune(".,!?:;", c) {
			if i > 0 && i+1 < len(lettres) && unicode.IsDigit(lettres[i-1]) && unicode.IsDigit(lettres[i+1]) {
				motActuel += string(c)
			} else {
				finirMot()
				mots = append(mots, string(c))
			}

		} else if c == '\'' {
			if i > 0 && i+1 < len(lettres) && unicode.IsLetter(lettres[i-1]) && unicode.IsLetter(lettres[i+1]) {
				motActuel += "'"
			} else {
				finirMot()
				mots = append(mots, "'")
			}

		} else {
			motActuel += string(c)
		}
	}
	finirMot()
	return mots
}

func trouverParenthese(lettres []rune, debut int) int {
	for j := debut + 1; j < len(lettres); j++ {
		if lettres[j] == ')' {
			return j
		}
	}
	return -1
}

func estPonctuation(mot string) bool {
	return len(mot) == 1 && strings.Contains(".,!?:;", mot)
}

func estGuillemet(mot string) bool {
	return mot == "'"
}

func estUnVraiMot(mot string) bool {
	return !estPonctuation(mot) && !estGuillemet(mot)
}

func lireBalise(texte string) (string, int, bool) {
	if len(texte) < 3 || texte[0] != '(' || texte[len(texte)-1] != ')' {
		return "", 0, false
	}

	interieur := texte[1 : len(texte)-1]
	parties := strings.Split(interieur, ",")

	nom := strings.TrimSpace(parties[0])
	if nom != "hex" && nom != "bin" && nom != "up" && nom != "low" && nom != "cap" {
		return "", 0, false
	}

	nombre := 1
	if len(parties) == 2 {
		n, err := strconv.Atoi(strings.TrimSpace(parties[1]))
		if err != nil || n < 1 {
			return "", 0, false
		}
		nombre = n
	} else if len(parties) > 2 {
		return "", 0, false
	}
	return nom, nombre, true
}

func estBalise(texte string) bool {
	_, _, ok := lireBalise(texte)
	return ok
}

func appliquerBalises(mots []string) []string {
	var resultat []string

	for _, mot := range mots {
		nom, nombre, ok := lireBalise(mot)
		if !ok {
			resultat = append(resultat, mot)
		} else {
			modifierMotsPrecedents(resultat, nom, nombre)
		}
	}
	return resultat
}

func modifierMotsPrecedents(mots []string, nom string, nombre int) {
	if nom == "hex" || nom == "bin" {
		base := 16
		if nom == "bin" {
			base = 2
		}
		for i := len(mots) - 1; i >= 0; i-- {
			if estUnVraiMot(mots[i]) {
				valeur, err := strconv.ParseInt(mots[i], base, 64)
				if err == nil {
					mots[i] = strconv.FormatInt(valeur, 10)
				}
				return
			}
		}
		return
	}

	compte := 0
	for i := len(mots) - 1; i >= 0 && compte < nombre; i-- {
		if !estUnVraiMot(mots[i]) {
			continue
		}
		if nom == "up" {
			mots[i] = strings.ToUpper(mots[i])
		} else if nom == "low" {
			mots[i] = strings.ToLower(mots[i])
		} else if nom == "cap" {
			mots[i] = mettreMajuscule(mots[i])
		}
		compte++
	}
}

func mettreMajuscule(mot string) string {
	lettres := []rune(strings.ToLower(mot))
	lettres[0] = unicode.ToUpper(lettres[0])
	return string(lettres)
}

func corrigerAAn(mots []string) {
	for i := 0; i < len(mots); i++ {
		if mots[i] != "a" && mots[i] != "A" {
			continue
		}
		j := i + 1
		for j < len(mots) && estGuillemet(mots[j]) {
			j++
		}
		if j >= len(mots) || !estUnVraiMot(mots[j]) {
			continue
		}

		premiereLettre := unicode.ToLower([]rune(mots[j])[0])
		if strings.ContainsRune("aeiouh", premiereLettre) {
			mots[i] = mots[i] + "n" // a -> an, A -> An
		}
	}
}

func recomposer(mots []string) string {
	resultat := ""
	pasDespace := true
	guillemetOuvert := false

	for _, mot := range mots {
		if estPonctuation(mot) {
			resultat += mot
			pasDespace = false

		} else if estGuillemet(mot) {
			if !guillemetOuvert {
				if !pasDespace {
					resultat += " "
				}
				resultat += mot
				guillemetOuvert = true
				pasDespace = true
			} else {
				resultat += mot
				guillemetOuvert = false
				pasDespace = false
			}

		} else {

			if !pasDespace {
				resultat += " "
			}
			resultat += mot
			pasDespace = false
		}
	}
	return resultat
}
