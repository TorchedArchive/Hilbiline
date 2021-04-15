// Hilbiline is a readline-inspired line editor made in pure Go

package hilbiline

import (
	"bufio"
	"os"

	_ "github.com/mattn/go-runewidth" // we'll need later
)

const (
	// Keycodes
	KeyNull      = 0
	KeyCtrlA     = 1
	KeyCtrlB     = 2
	KeyCtrlC     = 3
	KeyCtrlD     = 4
	KeyCtrlE     = 5
	KeyCtrlF     = 6
	KeyCtrlH     = 8
	KeyTab       = 9
	KeyCtrlK     = 11
	KeyCtrlL     = 12
	KeyEnter     = 13
	KeyCtrlN     = 14
	KeyCtrlP     = 16
	KeyCtrlT     = 20
	KeyCtrlU     = 21
	KeyCtrlW     = 23
	KeyEsc       = 27
	KeyBackspace = 127
)

var (
	maskedmode = 0
	mlmode     = 0
)

type HilbilineState struct {
	// io readers
	stdio  *bufio.Reader
	stdout *bufio.Reader

	// Readline buffer and prompt
	buf    []rune
	prompt string

	historyindex int

	pos    int // Current cursor position
	oldpos int // Previous cursor position
	cols   int // Num of terminal columns

	// Don't know if needed
	// Num of rows in mlmode
	maxrows int

	histState histState
}

func NewHilbilineState(prompt string) HilbilineState {
	return HilbilineState{
		stdio:  bufio.NewReader(os.Stdin),
		stdout: bufio.NewReader(os.Stdout),

		buf:    []rune{},
		prompt: prompt,

		// By default, does not have a file to write to. AddHistFile must be used for
		// persistent history
		histState: newHistState(),
	}
}

func (h HilbilineState) AddHistFile(path string) {
	// Open file with R/W perms or create it, perms are RWE for user only
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}

	h.histState.file = file

}
