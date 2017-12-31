package cache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	memcacheConnector := new(MemcacheConnector)

	mc := memcacheConnector.Connect(params)

	mc.Put("fo0", "bar", 1)

	// got, err := mc.Get("foo")

	// if err != nil {
	// 	panic(err)
	// }

	// if string(got.Value) != "my value" {
	// 	t.Error("Expected my value got", got)
	// }
}
