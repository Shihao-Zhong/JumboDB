package data_structure

import (
	"JumboDB/jumboDB-core/src/config"
	"JumboDB/jumboDB-core/src/protocol"
	"bufio"
	"errors"
	"github.com/bits-and-blooms/bloom/v3"
	"io"
	"log"
)

type Memtable struct {
	Data *SkipList
	BloomFilter *bloom.BloomFilter
	WALWriter *bufio.Writer
	Size int
}

func NewEmptyMemtable(config *config.TomlConfig) *Memtable {
	memtable := new(Memtable)
	memtable.Data = NewSkipList(config.Storage.SkipListLevel)
	memtable.BloomFilter = bloom.NewWithEstimates(uint(config.Storage.MemoryTableSize), config.Storage.BloomFilterFalsePositiveRate)
	memtable.WALWriter = openFileWithWriter(config.Storage.WALLocation)
	memtable.Size = 0
	return memtable
}

func GetMemtable(config *config.TomlConfig, currentIndex int) *Memtable {
	memtable := NewEmptyMemtable(config)
	if FileIsExist(config.Storage.WALLocation) {
		memtable.ReadFromWAL(config.Storage.WALLocation, currentIndex)
	}
	return memtable
}

func (i *Memtable) WriteWAL(opt *Operation) int {
	log.Printf("Start writing WAL for operation [%s]", opt.toString())
	n, err := i.WALWriter.Write(opt.OperationToJson())
	i.WALWriter.WriteString("\n")
	err = i.WALWriter.Flush()
	if err != nil {
		log.Printf("Error in write WAL for key [%s] with error [%s]\n", opt.Key, err)
	}
	i.Size += 1
	return n
}

func (i *Memtable) WriteDataToDisk(path string) int {
	return i.Data.toFile(path)
}

func (i *Memtable) Put(key string, value string) {
	operation := NewOperation(key, value, PUT)
	i.WriteWAL(operation)
	i.Data.Put(operation)
	i.BloomFilter.Add([]byte(key))
}

func (i *Memtable) Get(key string) (*Operation, error){
	if i.BloomFilter.Test([]byte(key)) {
		opt, err := i.Data.Get(key)
		if err != nil {
			return nil, errors.New("error in get from memtable")
		}
		return opt, nil
	}
	return nil, nil
}

func (i *Memtable) Apply(opt *Operation) {
	//i.WriteWAL(opt)
	i.Data.Put(opt)
	i.BloomFilter.Add([]byte(opt.Key))
}

// when restart db, readthe wal back from disk
func (i * Memtable) ReadFromWAL(walPath string, currentIndex int) {
	log.Printf("read wal from [%s]", walPath)
	scanner := openFileWithReader(walPath)
	// skip the wal that already in sstable
	for idx := 0; idx < currentIndex; idx++ {
		scanner.ReadString('\n')
	}
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			} else {
				panic(err)
			}
		}
		log.Printf("get data [%s]", line)
		opt := NewOperationFromJson(line)
		log.Printf("get opt [%s]", opt.toString())
		i.Apply(opt)
	}
}

func (i *Memtable) GetAll() []protocol.Payload {
	return i.Data.GetAll()
}

func (i *Memtable) Del(key string) {
	i.Data.Del(key)
}
