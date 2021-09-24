package service

import (
	"log"
	"../../../protocol"
)

func ParseAndExecuteRequest(bytePayload []byte) string {
	payload := DecodePayload(bytePayload)
	return executeRequest(payload)
}

func executeRequest(payload Payload) string {
	if payload.Operation = "get" {
		log.Println("get request issued")
	} else if payload.Operation = "put" {
		log.Println("put request issued")
	}
	return payload.Operation
}

