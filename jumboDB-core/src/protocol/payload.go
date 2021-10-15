package protocol

import (
	"encoding/json"
	"log"
)

type Payload struct {
	Key 	string `json:"key"`
	Value 	string `json:"value"`
}

func NewPayload(key string, value string) *Payload {
	payload := new(Payload)
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
