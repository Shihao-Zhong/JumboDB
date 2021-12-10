package data_structure

import (
	"encoding/json"
	"fmt"
)

type Operation struct {
	Key 	string `json:"key"`
	Value 	string `json:"value"`
	Operation string `json:"operation"`
	TransactionId int `json:"transactionId"`
}

const (
	PUT string = "PUT"
	DEL string = "DEL"
)

func NewOperation(key string, value string, operation string, trasactionId int) *Operation {
	op := new(Operation)
	op.Key = key
	op.Value = value
	op.Operation = operation
	op.TransactionId = trasactionId
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
	//elog.Printf("incoming json string [%s]\n", jsonString)
	var opt Operation
	err := json.Unmarshal([]byte(jsonString), &opt)
	if err != nil {
		panic(err)
	}
	return &opt
}

func (i *Operation) toString() string {
	if i != nil {
		return fmt.Sprintf("key = [%s], value = [%s], opt = [%s] transactionId = [%d]", i.Key, i.Value, i.Operation, i.TransactionId)
	}
	return fmt.Sprintf("key = [%v], value = [%v], opt = [%v]", nil, nil, nil)
}
