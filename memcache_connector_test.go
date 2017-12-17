package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	params := make(map[string]interface{})

	params["server 1"] = "localhost:11211"
	params["prefix"] = "golavel"

	memcacheConnector := new(MemcacheConnector)

	mc := memcacheConnector.Connect(params)

	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	got, err := mc.Get("foo")

	if err != nil {
		panic(err)
	}

	if string(got.Value) != "my value" {
		t.Error("Expected my value got", got)
	}
}
