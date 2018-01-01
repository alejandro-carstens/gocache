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

	mc.Put("foo", "bar", 1)

	got := mc.Get("foo")

	if got != "bar" {
		t.Error("Expected bar got", got)
	}
}
