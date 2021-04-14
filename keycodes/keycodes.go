package main

import (
	"fmt"
	"os"
	"bufio"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		input := ""
		quit := false
		for {
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()

			if err != nil {
				fmt.Println(err)
			}
			fmt.Print("\r\n", char, " // ", string(char))
			input = input + string(char)

			// print input on enter
			if char == 13 { fmt.Println(input); break }
			// quit on ctrl d
			if char == 4 { quit = true; break }
		}
		if quit == true { break }
	}
}
