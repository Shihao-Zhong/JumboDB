package service

import (
	"JumboDB/jumboDB-core/src/persistence"
	"JumboDB/jumboDB-core/src/protocol"
)

type Service struct {
	dataStorage persistence.Storage
}

func (i *Service)GetAllElements() []protocol.Payload {
	return i.dataStorage.GetAll()
}

func (i *Service) GetOneElement(key string) string {
	return i.dataStorage.Get(key)
}

func (i *Service) PutOneElement(key string, value string) {
	i.dataStorage.Put(key, value)
}

func (i *Service) DelOneElement(key string) {
	i.dataStorage.Del(key)
}

func (i *Service) Transaction(operations []map[string]string) error{
	return i.dataStorage.Transaction(operations)
}

func NewService(serviceType string) *Service {
	service := new(Service)
	switch serviceType {
	case "lsm": service.dataStorage = persistence.NewLSMTree()
	case "hashmap": service.dataStorage = persistence.NewInMemoryHashMap()
	}
	return service
}

