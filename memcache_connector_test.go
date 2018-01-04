package cache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	memcacheConnector := new(MemcacheConnector)

	memcacheStore := memcacheConnector.Connect(params)

	_, ok := memcacheStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", memcacheStore)
	}
}
