package encoder

import (
	"github.com/vmihailenco/msgpack/v5"
)

var _ Encoder = Msgpack{}

type Msgpack struct{}

// Encode implementation of the Encoder interface
func (Msgpack) Encode(item interface{}) ([]byte, error) {
	return msgpack.Marshal(item)
}

// Decode implementation of the Encoder interface
func (Msgpack) Decode(data []byte, dest interface{}) error {
	return msgpack.Unmarshal(data, dest)
}
