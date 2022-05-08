package encoder

type Encoder interface {
	Encode(item interface{}) ([]byte, error)
	Decode(data []byte, destination interface{}) error
}
