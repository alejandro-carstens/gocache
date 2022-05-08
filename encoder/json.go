package encoder

import (
	"encoding/json"
)

var _ Encoder = JSON{}

type JSON struct{}

// Encode implementation of the Encoder interface
func (JSON) Encode(item interface{}) ([]byte, error) {
	return json.Marshal(item)
}

// Decode implementation of the Encoder interface
func (JSON) Decode(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
