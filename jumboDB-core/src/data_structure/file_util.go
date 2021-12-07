package data_structure

import (
	"bufio"
	"log"
	"os"
)

func FileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func openFileWithWriter(filename string) *bufio.Writer {
	log.Printf("open writer in [%s]", filename)
	var f *os.File
	if FileIsExist(filename) {
		f, _ = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	} else {
		f, _ = os.Create(filename)
	}
	return bufio.NewWriter(f)
}

func openFileWithReader(filename string) *bufio.Reader {
	var f *os.File
	if FileIsExist(filename) {
		f, _ = os.OpenFile(filename, os.O_RDONLY, 0666)
	} else {
		f, _ = os.Create(filename)
	}
	return bufio.NewReader(f)
}

func openFileWithReadWriter(fileName string) *bufio.ReadWriter {
	return bufio.NewReadWriter(openFileWithReader(fileName), openFileWithWriter(fileName))
}