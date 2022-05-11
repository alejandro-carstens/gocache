package encoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var encoders = map[string]Encoder{
	"json":    JSON{},
	"msgpack": Msgpack{},
}

func TestEncode(t *testing.T) {
	type ex struct {
		Name     string
		LastName string
		Year     int
	}

	for name, e := range encoders {
		t.Run(name, func(t *testing.T) {
			b, err := e.Encode(&ex{
				Name:     "Ayrton",
				LastName: "Senna",
				Year:     1960,
			})
			require.NoError(t, err)

			var res ex
			require.NoError(t, e.Decode(b, &res))
			require.Equal(t, "Ayrton", res.Name)
			require.Equal(t, "Senna", res.LastName)
			require.Equal(t, 1960, res.Year)
		})
	}
}
