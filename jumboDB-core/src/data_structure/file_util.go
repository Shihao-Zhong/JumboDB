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

func OpenFileWithWriter(filename string) *bufio.Writer {
	log.Printf("open writer in [%s]", filename)
	var f *os.File
	if FileIsExist(filename) {
		f, _ = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	} else {
		f, _ = os.Create(filename)
	}
	return bufio.NewWriter(f)
}

func OpenFileWithReader(filename string) *bufio.Reader {
	var f *os.File
	if FileIsExist(filename) {
		f, _ = os.OpenFile(filename, os.O_RDONLY, 0666)
	} else {
		f, _ = os.Create(filename)
	}
	return bufio.NewReader(f)
}

func OpenFileWithReadWriter(fileName string) *bufio.ReadWriter {
	return bufio.NewReadWriter(OpenFileWithReader(fileName), OpenFileWithWriter(fileName))
}

func RemoveTables(tables []*Table) {
	for _, table := range tables {
		log.Printf("deleting file [%s]", table.FilePath)
		err := os.Remove(table.FilePath)
		if err != nil {
			log.Printf("Error in remove tables [%s]", err)
		}
	}
}