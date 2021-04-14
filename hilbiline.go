// Hilbiline is a readline-inspired line editor made in pure Go

package hilbiline

import (
	"io"
	"bufio"
	"os"
	"fmt"

	"golang.org/x/term"
	_ "github.com/mattn/go-runewidth" // we'll need later
)

const (
	KeyNull = 0
	KeyCtrlD = 4
	KeyTab = 9
	KeyEnter = 13
	KeyEsc = 27
	KeyBackspace = 127
)

type HilbilineState struct {
	buf []rune
	r *bufio.Reader
	prompt string
	pos int
}

func New(prompt string) HilbilineState {
	return HilbilineState{
		buf: []rune{},
		r: bufio.NewReader(os.Stdin),
		pos: 0,
		prompt: prompt,
	}
}

func (h *HilbilineState) Read() (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	h.buf = []rune{}
	h.pos = 0

	fmt.Print(h.prompt)
	for {
		char, _, err := h.r.ReadRune()
		if err != nil { return "", err }

		switch char {
		case KeyCtrlD:
			// return eof on ctrl d
			return "", io.EOF
		case KeyEnter:
			fmt.Print("\n\r")
			return string(h.buf), nil
		case KeyBackspace:
			h.backspace()
		default:
			// at default we assume is a printable char
			// so move the cursor (pos) and print the char
			h.pos++
			fmt.Print(string(char))
			h.buf = append(h.buf, char)
		}
	}

	// go complains if i dont have this
	return "", nil
}

// Sets the prompt
func (h *HilbilineState) SetPrompt(prompt string) {
	h.prompt = prompt
}

// backspace our text
// basically how it works is we move the cursor back, print
// empty space then go back again
// this doesnt work with characters with a length > 1
// ie japanese characters
func (h *HilbilineState) backspace() {
	if h.pos > 0 {
		h.buf = append(h.buf[:h.pos - 1], h.buf[h.pos:]...)
		h.pos--
		fmt.Printf("\u001b[1D \u001b[1D")
	}
}

