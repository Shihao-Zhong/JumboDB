package protocol

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

/*
Unit test for payload constructor
*/
func TestNewPayload(t *testing.T) {
	payload := NewPayload("a", "1")
	if payload.Key != "a" {
		t.Errorf("payload.Key excepted to be a, but %s got", payload.Key)
	} else if payload.Value != "1" {
		t.Errorf("payload.Value excepted to be 1, but %s got", payload.Value)
	}
}

/*
Unit test for EncodePayload
*/
func TestEncodePayload(t *testing.T) {
	payload := NewPayload("a", "1")
	expectJson := `{"Key":"a","Value":"1"}`
	payloadJson := string(EncodePayload(payload))
	if payloadJson != expectJson {
		t.Errorf("EncodePayload excepted to be %s, but %s got", expectJson, payloadJson)
	}
}

/*
Unit test for DecodePayload
*/
func TestDecodePayload(t *testing.T) {
	payload := NewPayload("a", "1")
	sourceJson := []byte(`{"Key":"a","Value":"1"}`)
	decodePayload := DecodePayload(sourceJson)
	if cmp.Equal(decodePayload, payload) {
		t.Errorf("DecodePayload excepted to be %s, but %s got", payload, decodePayload)
	}
}
