package hilbiline

import (
	"os"
)

const (
	hist_max_len = 1000
)

type histState struct {
	file    *os.File
	stifled bool
	histBuf histBuf
}

func newHistState() histState {
	return histState{
		stifled: false,
		histBuf: newHistBuf(),
	}
}

type histBuf struct {
	entries []string
	length  int
}

func newHistBuf() histBuf {
	return histBuf{
		entries: make([]string, 10),
		length:  0,
	}
}

// TODO
func (h histBuf) readFromFile(f *os.File) {}

func (h histBuf) writeToFile(f *os.File) {}
