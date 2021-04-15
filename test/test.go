package main

import (
	"hilbiline"
	"os"
)

func main() {
	homedir, _ := os.UserHomeDir()
	defaultconfpath := homedir + "/.hilbish-history"
	hl := hilbiline.NewHilbilineState("& ")
	hl.AddHistFile(defaultconfpath)
	// for {
	// str, e := hl.Read()
	// if e == io.EOF {
	// fmt.Println("hit ctrl d")
	// return
	// }
	//
	// if e != nil {
	// panic(e)
	// }
	//
	// fmt.Println(str)
	// }
}
