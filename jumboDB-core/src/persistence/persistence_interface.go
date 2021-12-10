package persistence

import "JumboDB/jumboDB-core/src/protocol"

type Storage interface {
	Get(key string) string
	Put(key string, value string) 
	GetAll() []protocol.Payload
	Del(key string)
	Transaction(operations []map[string]string) error
}
