package main

import (
	"fmt"
	"hilbiline"
	"io"
	"os"
)

func main() {
	homedir, _ := os.UserHomeDir()
	defaultconfpath := homedir + "/.hilbiline-history"

	hl := hilbiline.New("& ")
	hl.LoadHistory(defaultconfpath)

	for {
		str, e := hl.Read()

		if e == io.EOF {
			fmt.Println("hit ctrl d")
			return
		}

		if e != nil {
			panic(e)
		}

		fmt.Println(str)
	}
}
