package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Rosettea/Hilbiline"
)

func completion(line, cursinput string, pos int) []string {
	var suggestions []string
	items := []string{".git", ".gitignore", "hilbiline.go", "history.go"}
	for i := range items {
		if strings.HasPrefix(items[i], cursinput) {
			suggestions = append(suggestions, items[i][pos:])
		}
	}
	return suggestions
}

func main() {
	homedir, _ := os.UserHomeDir()
	defaultconfpath := homedir + "/.hilbiline-history"

	hl := hilbiline.New()
	hl.SetCompletionCallback(completion)
	hl.LoadHistory(defaultconfpath)

	for {
		str, e := hl.Read("\033[32m&\033[0m ")

		if e == io.EOF {
			fmt.Println("hit ctrl d")
			return
		}

		if e != nil {
			panic(e)
		}
		if str == "" { continue }

		fmt.Println(str)
	}
}
