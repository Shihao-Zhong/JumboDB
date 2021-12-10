package protocol

import (
	"encoding/json"
	"log"
)

type Payload struct {
	Key 	string `json:"key"`
	Value 	string `json:"value"`
	TransactionId int `json:"transactionId"`
}

func NewPayload(key string, value string, transactionId int) *Payload {
	payload := new(Payload)
	payload.Key = key
	payload.Value = value 
	payload.TransactionId = transactionId
	return payload
}

func EncodePayload(payload *Payload) []byte {
	b, err := json.Marshal(payload)
    if err != nil {
        log.Println(err)
    } 
	return b
}

func DecodePayload(payloadJson []byte) Payload {
	var payload Payload
	err := json.Unmarshal(payloadJson, &payload)
	if err != nil {
		log.Println(err)
	} 
	return payload
}

func DistinctMerge(payload1 []Payload, payload2 []Payload) []Payload {
	distinctMap := make(map[string]bool)
	for _, payload := range payload1 {
		distinctMap[payload.Key] = true
	}
	for _, payload := range payload2 {
		if _, ok := distinctMap[payload.Key]; ok {
			continue
		} else {
			payload1 = append(payload1, payload)
		}
	}
	return payload1
}