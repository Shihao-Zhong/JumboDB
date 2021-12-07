package data_structure

import (
	"JumboDB/jumboDB-core/src/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"io"
	"io/ioutil"
	"log"
)

// LevelSSTable
/*
Should persistent like the following config
{
    TotalLevel: int,
	Files: [[] string] string, each level should have a list of iles
	Sequences, the current naming sequence number
	Path, the default SSTable config location
}

*/
type LevelSSTable struct {
	TotalLevel int   `json:"totalLevel"`
	Data [][] *Table `json:"data"`
	Sequence int     `json:"sequence"`
	Path     string `json:"-"`
	
}

type Table struct {
	FilePath string `json:"filePath"`
	BloomFilterIndex *bloom.BloomFilter `json:"bloomFilterIndex"`
}

func (i *LevelSSTable) NewFileName(level int) string {
	seq := i.Sequence
	i.Sequence += 1
	return fmt.Sprintf("[%d]-[%d].jumbo", level, seq)
}

func NewEmptyLevelSSTable(path string) *LevelSSTable {
	levelSSTable := new(LevelSSTable)
	levelSSTable.TotalLevel = 0
	levelSSTable.Data = make([][]*Table, 1)
	levelSSTable.Sequence = 0
	levelSSTable.Path = path
	levelSSTable.ConfigToJsonFile()
	return levelSSTable
}

func GetLevelSSTable(path string) *LevelSSTable {
	if FileIsExist(path) {
		return ReadSSTableFromJsonFile(path)
	} else {
		return NewEmptyLevelSSTable(path)
	}
}

func (i *LevelSSTable) ConfigToJsonFile() {
	b, jsonMarshError := json.Marshal(i)
	if jsonMarshError != nil {
		panic(jsonMarshError)
	}

	ioError := ioutil.WriteFile(i.Path, b, 0755)
	if ioError != nil {
		log.Printf("SSTable ConfigToJsonFile write fail in location [%s]\n", i.Path)
		panic(ioError)
	}
}

func ReadSSTableFromJsonFile(path string) *LevelSSTable {
	var levelSSTable *LevelSSTable
	content ,ioError := ioutil.ReadFile(path)
	if ioError != nil {
		panic(ioError)
	}
	json.Unmarshal(content, &levelSSTable)
	levelSSTable.Path = path
	return levelSSTable
}

func (i *LevelSSTable) AddNewMemtable(memtable *Memtable) int {
	fileName := i.NewFileName(0)
	i.Data[0] = append([]*Table {NewTable(fileName, memtable.BloomFilter)}, i.Data[0]...)
	i.ConfigToJsonFile()
	return memtable.WriteDataToDisk(fileName)
}

func (i *LevelSSTable) majorCompaction(level int) {
	// need to avoid
	// first need to recrete file


}

func NewTable(filePath string, bloomFilterIndex *bloom.BloomFilter) *Table {
	levelTable := new(Table)
	levelTable.FilePath = filePath
	levelTable.BloomFilterIndex = bloomFilterIndex
	return levelTable
}

func (i *LevelSSTable) Get(key string) (*Operation, error) {
	for _, level := range i.Data {
		for _, table := range level {
			if table.BloomFilterIndex.Test([]byte(key)) {
				value, err := table.GetDataFromFile(key)
				if err != nil {
					continue
				} else {
					return value, nil
				}
			}
		}
	}
	return nil, errors.New("element not exist")
}

func (i *LevelSSTable) GetAll() []protocol.Payload {
	distinct := make(map[string]bool)
	var result []protocol.Payload
	for _, level := range i.Data {
		for _, table := range level {
			tableValue, err := table.GetAllDataFromFile()
			if err != nil {
				continue
			} else {
				for _, row := range tableValue {
					if _, ok := distinct[row.Key]; ok {
						continue
					} else {
						if row.Operation != DEL {
							result = append(result, *protocol.NewPayload(row.Key, row.Value))
						}
						distinct[row.Key] = true
					}
				}
			}
		}
	}
	return result
}

func (i *Table) GetDataFromFile(key string)  (*Operation, error) {
	scanner := openFileWithReader(i.FilePath)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("element not exist")
			} else {
				return nil, errors.New("error in read file")
			}
		}
		opt := NewOperationFromJson(line)
		if opt.Key == key {
			return opt, nil
		}
	}
}

func (i *Table) GetAllDataFromFile() ([]*Operation ,error) {
	var result []*Operation
	log.Printf("Getting data from file [%s]", i.FilePath)
	scanner := openFileWithReader(i.FilePath)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return result, nil
			} else {
				return nil, errors.New("error in read file")
			}
		}
		opt := NewOperationFromJson(line)
		result = append(result, opt)
	}
}



