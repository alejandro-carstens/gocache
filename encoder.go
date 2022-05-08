package gocache

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ Encoder = JSON{}
	_ Encoder = Msgpack{}
)

type (
	Encoder interface {
		Encode(item interface{}) ([]byte, error)
		Decode(data []byte, destination interface{}) error
	}
	JSON    struct{}
	Msgpack struct{}
)

// Encode implementation of the Encoder interface
func (Msgpack) Encode(item interface{}) ([]byte, error) {
	return msgpack.Marshal(item)
}

// Decode implementation of the Encoder interface
func (Msgpack) Decode(data []byte, dest interface{}) error {
	return msgpack.Unmarshal(data, dest)
}

// Encode implementation of the Encoder interface
func (JSON) Encode(item interface{}) ([]byte, error) {
	return json.Marshal(item)
}

// Decode implementation of the Encoder interface
func (JSON) Decode(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
