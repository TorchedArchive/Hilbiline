package hilbiline

import (
	"os"
)

var (
	HistoryMaxLen = 1000
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
	// Maybe have map[string]int?
	entries []string
	length  int
}

func newHistBuf() histBuf {
	return histBuf{
		entries: []string{},
		length:  0,
	}
}

// TODO
func (h *histBuf) readFromFile(f *os.File) {}

// TODO: Currently overwrites history
func (h *histBuf) writeToFile(f *os.File) {
	for _, v := range h.entries {
		f.WriteString(v)
		f.WriteString("\n")
	}
}

func (h *histBuf) addEntry(s string) {
	h.entries = append(h.entries, s)
}
