package persistence

import (
	"JumboDB/jumboDB-core/src/data_structure"
	"JumboDB/jumboDB-core/src/protocol"
	"log"
)

type InMemoryHashMap struct {
	dataStorage map[string]string
}

func(i InMemoryHashMap) Get(key string) string {
	value, ok := i.dataStorage[key]
	if (ok) {
		return value 
	} else {
		log.Printf("The value with key [%s] is not exist", key)
	}
	return ""
}

func(i InMemoryHashMap) Del(key string) {
	_, ok := i.dataStorage[key]
	if (ok) {
		delete(i.dataStorage, key)
	} else {
		log.Printf("The value with key [%s] is not exist", key)
	}
	return 
}

func(i InMemoryHashMap) Put(key string, value string) {
	i.dataStorage[key] = value
	return
}

func(i InMemoryHashMap) GetAll() []protocol.Payload {
	var data []protocol.Payload
	for k,v := range i.dataStorage {
		data = append(data, *protocol.NewPayload(k, v, 0))
	}
	return data
}

func NewInMemoryHashMap() *InMemoryHashMap {
	storage := new(InMemoryHashMap)
	storage.dataStorage = make(map[string]string)

	return storage
}
func (i InMemoryHashMap) Transaction(operations []map[string]string) error{
	for _, opt := range operations {
		if opt["operation"] == data_structure.PUT {
			i.dataStorage[opt["key"]] = opt["value"]
		} else {
			i.dataStorage[opt["key"]] = ""
		}
	}
	return nil
}
