package service

import (
	"JumboDB/jumboDB-core/src/persistence"
	"JumboDB/jumboDB-core/src/protocol"
)

var dataStorage = persistence.NewInMemoryHashMap()


func GetAllElements() []protocol.Payload {
	return dataStorage.GetAll()
}

func GetOneElement(key string) string {
	return dataStorage.Get(key)
}

func PutOneElement(key string, value string) {
	dataStorage.Put(key, value)
}

func DelOneElement(key string) {
	dataStorage.Del(key)
}

