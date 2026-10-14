package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type kind int

const (
	word kind = iota
	punct
	quote
	tag
)

type token struct {
	kind kind
	text string
}

var (
	tokenRe = regexp.MustCompile(`\((?:hex|bin|up|low|cap)(?:,\s*\d+)?\)|\d+(?:[.,]\d+)+|[.,!?:;]+|'|[^\s.,!?:;']+(?:'[^\s.,!?:;']+)*`)
	tagRe   = regexp.MustCompile(`^\((hex|bin|up|low|cap)(?:,\s*(\d+))?\)$`)
)

func Process(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = processLine(line)
	}
	return strings.Join(lines, "\n")
}

func processLine(line string) string {
	tokens := tokenize(line)
	tokens = applyTags(tokens)
	fixArticles(tokens)
	return render(tokens)
}

func tokenize(line string) []token {
	var tokens []token
	for _, m := range tokenRe.FindAllString(line, -1) {
		switch {
		case tagRe.MatchString(m):
			tokens = append(tokens, token{tag, m})
		case m == "'":
			tokens = append(tokens, token{quote, m})
		case strings.ContainsRune(".,!?:;", rune(m[0])):
			tokens = append(tokens, token{punct, m})
		default:
			tokens = append(tokens, token{word, m})
		}
	}
	return tokens
}

func applyTags(tokens []token) []token {
	var out []token
	for _, t := range tokens {
		if t.kind != tag {
			out = append(out, t)
			continue
		}
		sub := tagRe.FindStringSubmatch(t.text)
		name, n := sub[1], 1
		if sub[2] != "" {
			n, _ = strconv.Atoi(sub[2])
		}
		applyTag(out, name, n)
	}
	return out
}

func applyTag(tokens []token, name string, n int) {
	count := 0
	for i := len(tokens) - 1; i >= 0 && count < n; i-- {
		if tokens[i].kind != word {
			continue
		}
		switch name {
		case "hex", "bin":
			base := 16
			if name == "bin" {
				base = 2
			}
			if v, err := strconv.ParseInt(tokens[i].text, base, 64); err == nil {
				tokens[i].text = strconv.FormatInt(v, 10)
			}
			return
		case "up":
			tokens[i].text = strings.ToUpper(tokens[i].text)
		case "low":
			tokens[i].text = strings.ToLower(tokens[i].text)
		case "cap":
			tokens[i].text = capitalize(tokens[i].text)
		}
		count++
	}
}

func capitalize(s string) string {
	r := []rune(strings.ToLower(s))
	if len(r) > 0 {
		r[0] = unicode.ToUpper(r[0])
	}
	return string(r)
}

func fixArticles(tokens []token) {
	for i, t := range tokens {
		if t.kind != word || (t.text != "a" && t.text != "A") {
			continue
		}
		j := i + 1
		for j < len(tokens) && tokens[j].kind == quote {
			j++
		}
		if j >= len(tokens) || tokens[j].kind != word {
			continue
		}
		first := unicode.ToLower([]rune(tokens[j].text)[0])
		if strings.ContainsRune("aeiouh", first) {
			tokens[i].text += "n"
		}
	}
}

func render(tokens []token) string {
	var sb strings.Builder
	inQuote, noSpace := false, true
	for _, t := range tokens {
		switch t.kind {
		case word:
			if !noSpace {
				sb.WriteByte(' ')
			}
			sb.WriteString(t.text)
			noSpace = false
		case punct:
			sb.WriteString(t.text)
			noSpace = false
		case quote:
			if !inQuote {
				if !noSpace {
					sb.WriteByte(' ')
				}
				sb.WriteString(t.text)
				inQuote, noSpace = true, true
			} else {
				sb.WriteString(t.text)
				inQuote, noSpace = false, false
			}
		}
	}
	return sb.String()
}
