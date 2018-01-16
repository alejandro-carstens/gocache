package cache

import (
	"errors"
)

// ArrayConnector is a representation of the array store connector
type ArrayConnector struct{}

// Connect is responsible for connecting with the caching store
func (ac *ArrayConnector) Connect(params map[string]interface{}) (StoreInterface, error) {
	params, err := ac.validate(params)

	if err != nil {
		return &ArrayStore{}, err
	}

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &ArrayStore{
		Client: make(map[string]interface{}),
		Prefix: prefix,
	}, nil
}

func (ac *ArrayConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("You need to specify a caching prefix.")
	}

	return params, nil
}
