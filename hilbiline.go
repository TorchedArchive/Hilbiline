// Hilbiline is a readline-inspired line editor made in pure Go

package hilbiline

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

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
	KeyCtrlJ     = 10
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

		buf:    []rune{},
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

	h.buf = []rune{}
	h.pos = 0

	for {
		// TODO: custom read so we can handle esc properly
		char, _, err := h.stdio.ReadRune()
		if err != nil {
			return "", err
		}

		switch char {
		case KeyCtrlD:
			// End session on CtrlD
			h.destroy()
			return "", io.EOF
		case KeyCtrlC:
			fmt.Printf("\r\n")
			return "", nil
		// Vertical feed
		case KeyCtrlJ:
			return string(h.buf), nil
		case KeyCtrlL:
			h.ClearScreen()
		case KeyCtrlU:
			// Delete whole line
			h.buf = []rune{}
			h.pos = 0

			h.refreshLine()
		// case KeyCtrlK: remove reset of word
		// case KeyCtrlF: move forward one character
		// case KeyCtrlB: move back one character
		// case KeyCtrlT: swap buf[-1] and buf[-2]
		// case KeyCtrlH: delete last rune
		// case KeyCtrlN: go forward in history
		// case KeyCtrlP: go back in history
		// case KeyCtrlW: delete to previous word
		// case KeyEsc: handle esc codes (cursor up etc)
		case KeyEnter:
			fmt.Print("\n\r")
			h.histState.histBuf.addEntry(string(h.buf))
			return string(h.buf), nil
		case KeyBackspace:
			h.editBackspace()
		case KeyTab:
			// TODO: tab completion
			// Just making it a no-op so it wont print the tab
		default:
			h.editInsert(char)
		}
	}

}

func (h *HilbilineState) LoadHistory(path string) error {
	// Open file with R/W perms or create it,
	// perms are RW for user only
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	h.histState.file = file
	return nil
}

func (h *HilbilineState) PrintPrompt() {
	fmt.Print(h.prompt)
}

func (h *HilbilineState) SetPrompt(prompt string) {
	h.prompt = prompt
}

func (h *HilbilineState) ClearScreen() {
	fmt.Print("\x1b[H\x1b[2J")
	h.PrintPrompt()
}

func (h *HilbilineState) editInsert(c rune) {
	h.pos++
	h.buf = append(h.buf, c)

	if !mlmode {
		fmt.Print(string(c))
	} else {
		fmt.Print("*")
	}
}

func (h *HilbilineState) editBackspace() {
	if h.pos > 0 {
		_, length := utf8.DecodeLastRuneInString(string(h.buf))
		h.buf = append(h.buf[:h.pos-1], h.buf[h.pos:]...)
		h.pos--

		// This is atrocious
		// For testing, ん has length 3, English chars have 1
		// Without this if check ん will loop 4 times instead of 3
		if length == 1 {
			fmt.Print("\b")
		} else {
			for i := 1; i < length; i++ {
				// Backspace code
				fmt.Print("\b")
			}
		}

		// Clear to end
		fmt.Print("\033[K")
	}
}

func (h *HilbilineState) refreshLine() (*term.State, error) {
	// TODO: dont do this here; make a separate function
	// refreshLine should do as the name suggests, only refresh
	// the line
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	return oldState, err
}

func (h *HilbilineState) destroy() {
	h.histState.histBuf.writeToFile(h.histState.file)
}
