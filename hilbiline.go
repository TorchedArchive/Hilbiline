// Hilbiline is a readline-inspired line editor made in pure Go

package hilbiline

import (
	"bufio"
	"os"

	_ "github.com/mattn/go-runewidth" // we'll need later
)

const (
	KeyNull      = 0
	KeyCtrlD     = 4
	KeyTab       = 9
	KeyEnter     = 13
	KeyEsc       = 27
	KeyBackspace = 127
)

type HilbilineState struct {
	// Line state
	point int
	end   int
	mark  int
	// buflen int (not needed?)
	buf    []rune
	prompt string

	// Global state
	rlstate int
	done    int
	// keymap Keymap

	// Input state
	// kseq *rune ???
	// kseqlen int ???

	// pendingin int ???

	r *bufio.Reader
}

func New(prompt string) HilbilineState {
	return HilbilineState{
		point:  0,
		end:    0,
		mark:   0,
		buf:    []rune{},
		prompt: prompt,

		rlstate: 0,
		done:    0,

		r: bufio.NewReader(os.Stdin),
	}
}
