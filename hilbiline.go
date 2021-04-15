// Hilbiline is a readline-inspired line editor made in pure Go

package hilbiline

import (
	"bufio"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-runewidth" // we'll need later
	"golang.org/x/term"
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
	maskedmode = false
	mlmode     = false
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

func New(prompt string) HilbilineState {
	return HilbilineState{
		stdio:  bufio.NewReader(os.Stdin),
		stdout: bufio.NewReader(os.Stdout),

		// Preallocate to avoid reallocation later
		buf:    make([]rune, 80),
		prompt: prompt,

		// By default, does not have a file to write to.
		// AddHistFile must be used for persistent history
		histState: newHistState(),
	}
}

func (h *HilbilineState) Read() (string, error) {
	fmt.Print(h.prompt)

	oldState, _ := h.refreshLine()
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	h.buf = make([]rune, 80)
	h.pos = KeyNull

	for {
		// TODO: custom read so we can handle esc properly
		char, _, err := h.stdio.ReadRune()
		if err != nil {
			return "", err
		}

		switch char {
		// Apparently CtrlT, CtrlB, CtrlF all do stuff
		// but like no one cares about it
		case KeyCtrlD:
			// End session on CtrlD
			return "", io.EOF
		case KeyCtrlC:
			return "", nil
		case KeyCtrlL:
			h.ClearScreen()
		case KeyCtrlU:
			// Delete whole line
			h.buf[0] = KeyNull
			h.pos = 0
			h.refreshLine()
		// case KeyCtrlN: go forward in history
		// case KeyCtrlP: go back in history
		// case KeyEsc: handle esc codes (cursor up etc)
		case KeyEnter:
			fmt.Print("\n\r")
			return string(h.buf), nil
		case KeyBackspace:
			h.editBackspace()
		default:
			h.pos++
			h.editInsert(char)
		}
	}

}

func (h HilbilineState) LoadHistory(path string) error {
	// Open file with R/W perms or create it,
	// perms are RWE for user only
	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0700)
	if err != nil {
		return err
	}

	h.histState.file = file
	return nil
}

func (h HilbilineState) PrintPrompt() {
	fmt.Print(h.prompt)
}

func (h HilbilineState) ClearScreen() {
	fmt.Print("\x1b[H\x1b[2J")
	h.PrintPrompt()
}

func (h HilbilineState) editInsert(c rune) {
	h.buf[h.pos] = c
	h.pos++
	h.buf[h.pos] = KeyNull

	if !mlmode {
		fmt.Print(string(c))
	} else {
		fmt.Print("*")
	}
}

func (h HilbilineState) editBackspace() {
	h.pos--
	h.buf[h.pos] = KeyNull
	h.refreshLine()
}

func (h *HilbilineState) refreshLine() (*term.State, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	return oldState, err
}
