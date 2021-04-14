package main

import (
	"io"
	"hilbiline"
	"fmt"
)

func main() {
	hl := hilbiline.New("& ")
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
