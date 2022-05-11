package encoder

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ Encoder = JSON{}
	_ Encoder = Msgpack{}
)

type (
	// Encoder represents an interface that exposes functionality to encode and decode non-numeric or
	// boolean cache entries
	Encoder interface {
		// Encode returns the encoded/marshaled version of item
		Encode(item interface{}) ([]byte, error)
		// Decode decodes the MessagePack-encoded data and stores the result
		// in the value pointed to by v.
		Decode(data []byte, destination interface{}) error
	}
	// JSON is an Encoder implementation for the encoding/json package
	JSON struct{}
	// Msgpack is an Encoder implementation for the msgpack package. To learn
	// more about msgpack please see: https://msgpack.uptrace.dev
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
