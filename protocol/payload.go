package protocol

import (
	"encoding/json"
	"log"
)

type Payload struct {
	Operation      string
	Key 	string
	Value 	string
}

func NewPayload(operation string, key string, value string) *Payload {
	payload := new(Payload)
	payload.Operation = operation
	payload.Key = key
	payload.Value = value 
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
