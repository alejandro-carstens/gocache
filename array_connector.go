package cache

import (
	"errors"
)

type ArrayConnector struct{}

func (this *ArrayConnector) Connect(params map[string]interface{}) (StoreInterface, error) {
	params, err := this.validate(params)

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

func (this *ArrayConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("You need to specify a caching prefix.")
	}

	return params, nil
}
