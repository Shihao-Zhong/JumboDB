package persistence

import (
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
		log.Println("The value with key [%s] is not exist", key)
	}
	return ""
}

func(i InMemoryHashMap) Del(key string) {
	_, ok := i.dataStorage[key]
	if (ok) {
		delete(i.dataStorage, key)
	} else {
		log.Println("The value with key [%s] is not exist", key)
	}
	return 
}

func(i InMemoryHashMap) Put(key string, value string) {
	i.dataStorage[key] = value
	return
}

func(i InMemoryHashMap) GetAll() []protocol.Payload{
	var data []protocol.Payload
	for k,v := range i.dataStorage {
		data = append(data, *protocol.NewPayload(k, v))
	}
	return data
}

func NewInMemoryHashMap() *InMemoryHashMap {
	storage := new(InMemoryHashMap)
	storage.dataStorage = make(map[string]string)

	return storage
}

