package persistence

import (
	"JumboDB/jumboDB-core/src/config"
	"JumboDB/jumboDB-core/src/data_structure"
	"JumboDB/jumboDB-core/src/protocol"
	"encoding/json"
	"io/ioutil"
	"log"
)

type LSMTree struct {
	Memtable            *data_structure.Memtable `json:"-"`
	ImmutableMemtable 	*data_structure.Memtable `json:"-"`
	Config              *config.TomlConfig `json:"-"`
	LSSTable            *data_structure.LevelSSTable `json:"-"`
	Index *LSMTreeIndex
}

type LSMTreeIndex struct {
	CurrentIndex int `json:"currentIndex"`
}

// NewEmptyLSMTree Create a pure new lsm tree, when database first time running.
func NewLSMTree() *LSMTree {
	lsm := new(LSMTree)
	lsm.Config = config.GetConfig()
	lsm.Index = GetLSMTreeIndex(lsm.Config.Storage.LSMIndexLocation)
	lsm.Memtable = data_structure.GetMemtable(lsm.Config, lsm.Index.CurrentIndex	)
	lsm.LSSTable = data_structure.GetLevelSSTable(lsm.Config.Storage.SSTableIndexLocation)

	return lsm
}

func (i *LSMTreeIndex) ConfigToJsonFile(path string) {
	b, jsonMarshError := json.Marshal(i)
	if jsonMarshError != nil {
		panic(jsonMarshError)
	}

	ioError := ioutil.WriteFile(path, b, 0755)
	if ioError != nil {
		log.Printf("LsmTree ConfigToJsonFile write fail in location [%s]\n", path)
		panic(ioError)
	}
}

func NewLSMTreeIndex () *LSMTreeIndex {
	lsmIndex := new(LSMTreeIndex)
	lsmIndex.CurrentIndex = 0
	return lsmIndex
}

func ReadLSMTreeIndexFromJsonFile(path string) *LSMTreeIndex {
	var lsmIndex *LSMTreeIndex
	content ,ioError := ioutil.ReadFile(path)
	if ioError != nil {
		panic(ioError)
	}
	json.Unmarshal(content, &lsmIndex)
	return lsmIndex
}

func GetLSMTreeIndex(path string) *LSMTreeIndex{
	if data_structure.FileIsExist(path) {
		return ReadLSMTreeIndexFromJsonFile(path)
	}
	return NewLSMTreeIndex()
}

func (i *LSMTree) Put(key string, value string) {
	// just need to find first available memtable to write. writing thread will help to clear all things.
	i.Memtable.Put(key, value)
	log.Printf("current memrable size = [%d]", i.Memtable.Size)
	if i.Memtable.Size > i.Config.Storage.MemoryTableSize {
		go i.SwapMemtable()
	}
}

func (i *LSMTree) SwapMemtable() {
	// if immutable is now write to disk or some other thread already swap
	if i.ImmutableMemtable != nil {
		return
	}
	// because of pointer change, the existing call should be good
	i.ImmutableMemtable = i.Memtable
	i.Memtable = data_structure.NewEmptyMemtable(i.Config)
	writeLine := i.LSSTable.AddNewMemtable(i.ImmutableMemtable)
	i.ImmutableMemtable = nil
	i.Index.CurrentIndex += writeLine
	i.Index.ConfigToJsonFile(i.Config.Storage.LSMIndexLocation)
}

func (i *LSMTree) Get(key string) string{
	value, err := i.Memtable.Get(key)
	if err != nil {
		log.Printf("Error in lsmtree get for key [%s] with error [%s]\n", key, err)
	}
	if value == nil {
		// immutable memtable is not release yet
		if i.ImmutableMemtable != nil {
			value, err = i.ImmutableMemtable.Get(key)
			if err != nil {
				log.Printf("Error in lsmtree immutable memtable get for key [%s] with error [%s]\n", key, err)
			}
		}
		if value == nil {
			value, err = i.LSSTable.Get(key)
			if err != nil {
				log.Printf("Error in lsmtree sstable get for key [%s] with error [%s]\n", key, err)
			}
		}
	}
	// if find
	if value != nil {
		if value.Operation == data_structure.DEL {
			return ""
		} else {
			return value.Value
		}
	}
	// not find
	return ""
}

func (i *LSMTree) Del(key string) {
	i.Memtable.Del(key)
}

func (i *LSMTree) GetAll() []protocol.Payload {
	log.Printf("SStable data = [%s]", i.LSSTable.GetAll())
	return protocol.DistinctMerge(i.Memtable.GetAll(), i.LSSTable.GetAll())
}
