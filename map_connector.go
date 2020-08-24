package gocache

import (
	"errors"
)

// MapConnector is a representation of the array store connector
type MapConnector struct{}

// Connect is responsible for connecting with the caching store
func (ac *MapConnector) Connect(params map[string]interface{}) (Cache, error) {
	params, err := ac.validate(params)
	if err != nil {
		return nil, err
	}

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &MapStore{
		client: make(map[string]interface{}),
		prefix: prefix,
	}, nil
}

func (ac *MapConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("you need to specify a caching prefix")
	}

	return params, nil
}
