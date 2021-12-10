package data_structure

import (
	"JumboDB/jumboDB-core/src/config"
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
	isCompacting bool `json:"-"`
	
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
	levelSSTable.isCompacting = false
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

func (i *LevelSSTable) AddNewMemtable(memtable *Memtable, config *config.TomlConfig) int {
	fileName := i.NewFileName(0)
	i.Data[0] = append([]*Table {NewTable(fileName, memtable.BloomFilter)}, i.Data[0]...)
	if len(i.Data[0]) >= config.Storage.MajorCompactionFileSize {
		go i.MajorCompaction(config)
	}
	i.ConfigToJsonFile()
	return memtable.WriteDataToDisk(fileName)
}

func (i *LevelSSTable) MajorCompaction(config *config.TomlConfig) {
	log.Printf("Start major compaction")
	if i.isCompacting {
		return
	}
	// need to avoid
	// first need to recrete file

	for level, files := range i.Data {
		var fileNumber = len(files)
		if fileNumber < config.Storage.MajorCompactionFileSize {
			continue
		}
		var tableList = make([]*Table, fileNumber)
		tableList = files[len(files)-fileNumber:]
		var fileName = i.NewFileName(level+1)
		newTable := multiWayMerge(tableList, fileName, fileNumber, config)

		newLevel := []*Table{newTable}
		// if last level create new level else add to next level
		if level == len(i.Data) - 1 {
			i.Data = append(i.Data, newLevel)
		} else {
			i.Data[level+1] = append(newLevel, i.Data[level+1]...)
		}
		// remove all existing files in this level
		i.Data[level] = files[:len(files)-fileNumber]
		i.ConfigToJsonFile()
		go RemoveTables(tableList)
	}
	i.isCompacting = false
}

func NewTable(filePath string, bloomFilterIndex *bloom.BloomFilter) *Table {
	levelTable := new(Table)
	levelTable.FilePath = filePath
	levelTable.BloomFilterIndex = bloomFilterIndex
	return levelTable
}

func (i *LevelSSTable) Get(key string, transactionId int) (*Operation, error) {
	for _, level := range i.Data {
		for _, table := range level {
			if table.BloomFilterIndex.Test([]byte(key)) {
				value, err := table.GetDataFromFile(key, transactionId)
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
							result = append(result, *protocol.NewPayload(row.Key, row.Value, row.TransactionId))
						}
						distinct[row.Key] = true
					}
				}
			}
		}
	}
	return result
}

func (i *Table) GetDataFromFile(key string, transactionId int)  (*Operation, error) {
	scanner := OpenFileWithReader(i.FilePath)
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
		if opt.Key == key && opt.TransactionId < transactionId{
			return opt, nil
		}
	}
}

func (i *Table) GetAllDataFromFile() ([]*Operation ,error) {
	var result []*Operation
	log.Printf("Getting data from file [%s]", i.FilePath)
	scanner := OpenFileWithReader(i.FilePath)
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



