package persistence

import (
	"JumboDB/jumboDB-core/src/config"
	"JumboDB/jumboDB-core/src/data_structure"
	"JumboDB/jumboDB-core/src/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type LSMTree struct {
	Memtable            *data_structure.Memtable `json:"-"`
	ImmutableMemtable 	*data_structure.Memtable `json:"-"`
	Config              *config.TomlConfig `json:"-"`
	LSSTable            *data_structure.LevelSSTable `json:"-"`
	Index *LSMTreeIndex
	TransactionLocks map[string] int `json:"-"` // key = key value = transaction Id, do 2pl in w/w lock
}

type LSMTreeIndex struct {
	CurrentIndex int `json:"currentIndex"`
	CurrentTransactionIndex int `json:"currentTransactionIndex"`
	TransactionIndexLock sync.Mutex `json:"-"`
}

// NewEmptyLSMTree Create a pure new lsm tree, when database first time running.
func NewLSMTree() *LSMTree {
	lsm := new(LSMTree)
	lsm.Config = config.GetConfig()
	lsm.Index = GetLSMTreeIndex(lsm.Config.Storage.LSMIndexLocation)
	lsm.Memtable = data_structure.GetMemtable(lsm.Config, lsm.Index.CurrentIndex	)
	lsm.LSSTable = data_structure.GetLevelSSTable(lsm.Config.Storage.SSTableIndexLocation)
	lsm.TransactionLocks = make(map[string] int)
	return lsm
}

func (i *LSMTreeIndex) NextTransactionIndex() int {
	i.TransactionIndexLock.Lock()
	defer i.TransactionIndexLock.Unlock()

	transactionId := i.CurrentTransactionIndex
	i.CurrentTransactionIndex++
	return transactionId
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
	lsmIndex.CurrentTransactionIndex = 0
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
	transactionId := i.Index.NextTransactionIndex()
	i.putTransaction(key, value, transactionId)
}

func (i *LSMTree) putTransaction(key string, value string, transactionId int) {
	i.TransactionLocks[key] = transactionId
	i.Memtable.Put(key, value, transactionId)
	log.Printf("current memtable size = [%d] with max size [%d]", i.Memtable.Size, i.Config.Storage.MemoryTableSize)
	delete(i.TransactionLocks, key)
	if i.Memtable.Size >= i.Config.Storage.MemoryTableSize {
		go i.SwapMemtable()
	}
}

func (i *LSMTree) SwapMemtable() {
	log.Printf("Start swap memtable with size %d", i.Memtable.Size)
	// if immutable is now write to disk or some other thread already swap
	if i.ImmutableMemtable != nil {
		log.Printf("Other thread is using immutable memtable")
		return
	}
	// because of pointer change, the existing call should be good
	i.ImmutableMemtable = i.Memtable
	i.Memtable = data_structure.NewEmptyMemtable(i.Config)
	writeLine := i.LSSTable.AddNewMemtable(i.ImmutableMemtable, i.Config)
	i.ImmutableMemtable = nil
	i.Index.CurrentIndex += writeLine
	i.Index.ConfigToJsonFile(i.Config.Storage.LSMIndexLocation)
	log.Printf("Finish write memtable.")
}

func (i *LSMTree) Get(key string) string{
	transactionId := i.Index.NextTransactionIndex()
	value, err := i.Memtable.Get(key, transactionId)
	if err != nil {
		log.Printf("Error in lsmtree get for key [%s] with error [%s]\n", key, err)
	}

	if value == nil {
		// immutable memtable is not release yet
		log.Printf("Error in lsmtree get for key [%s] did not exist in memtable\n", key)
		if i.ImmutableMemtable != nil {
			value, err = i.ImmutableMemtable.Get(key, transactionId)
			if err != nil {
				log.Printf("Error in lsmtree immutable memtable get for key [%s] with error [%s]\n", key, err)
			}
		}
		if value == nil {
			value, err = i.LSSTable.Get(key, transactionId)
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
	transactionId := i.Index.NextTransactionIndex()
	i.delTransaction(key, transactionId)
}

func (i *LSMTree) delTransaction(key string, transactionId int) {
	i.Memtable.Del(key, transactionId)
	log.Printf("current memtable size = [%d] with max size [%d]", i.Memtable.Size, i.Config.Storage.MemoryTableSize)
	if i.Memtable.Size >= i.Config.Storage.MemoryTableSize {
		go i.SwapMemtable()
	}
}

func (i *LSMTree) GetAll() []protocol.Payload {
	log.Printf("SStable data = [%s]", i.LSSTable.GetAll())
	return protocol.DistinctMerge(i.Memtable.GetAll(), i.LSSTable.GetAll())
}

func (i *LSMTree) Transaction(operations []map[string]string) error {
	transactionId := i.Index.NextTransactionIndex()
	locksKey := make([]string, len(operations))
	// step one, lock all rows
	for idx, opt := range operations {
		// locking step
		if value, ok := i.TransactionLocks[opt["key"]]; ok {
			if value != transactionId {
				// if some row already locked by other transaction, revert
				// first release all locks
				for _, lock := range locksKey {
					delete(i.TransactionLocks, lock)
				}
				return errors.New(fmt.Sprintf("row [%s] has already be locked by transaction [%d]", opt["key"], value))
			}
		} else {
			// lock the row
			i.TransactionLocks[opt["key"]] = transactionId
		}
		locksKey[idx] = opt["key"]
	}
	// do transaction and release lock
	for idx, opt := range operations {
		// transaction step
		if opt["operation"] == data_structure.PUT {
			i.putTransaction(opt["key"], opt["value"], transactionId)
		} else {
			i.delTransaction(opt["key"], transactionId)
		}
		// release lock
		delete(i.TransactionLocks, locksKey[idx])
	}
	return nil
}

