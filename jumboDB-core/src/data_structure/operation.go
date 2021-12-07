package data_structure

import (
	"encoding/json"
	"fmt"
)

type Operation struct {
	Key 	string `json:"key"`
	Value 	string `json:"value"`
	Operation string `json:"operation"`
}

const (
	PUT string = "PUT"
	DEL string = "DEL"
)

func NewOperation(key string, value string, operation string) *Operation {
	op := new(Operation)
	op.Key = key
	op.Value = value
	op.Operation = operation
	return op
}

func (i *Operation) OperationToJson() []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return b
}

func NewOperationFromJson(jsonString string) *Operation {
	//log.Printf("incoming json string [%s]\n", jsonString)
	var opt Operation
	err := json.Unmarshal([]byte(jsonString), &opt)
	if err != nil {
		panic(err)
	}
	return &opt
}

func (i *Operation) toString() string {
	if i != nil {
		return fmt.Sprintf("key = [%s], value = [%s], opt = [%s]", i.Key, i.Value, i.Operation)
	}
	return fmt.Sprintf("key = [%v], value = [%v], opt = [%v]", nil, nil, nil)
}
