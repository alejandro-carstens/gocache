package gocache

import "errors"

// New new-ups an instance of Store
func New(config *Config) (Cache, error) {
	var connector cacheConnector
	if config.Map != nil {
		connector = new(mapConnector)
	}
	if config.Memcache != nil {
		connector = new(memcacheConnector)
	}
	if config.Redis != nil {
		connector = new(redisConnector)
	}
	if connector == nil {
		return nil, errors.New("invalid or empty config specified")
	}
	if err := connector.validate(config); err != nil {
		return nil, err
	}

	return connector.connect(config)
}
