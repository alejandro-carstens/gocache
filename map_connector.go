package gocache

type mapConnector struct{}

func (ac *mapConnector) connect(config *Config) (Cache, error) {
	return &MapStore{
		client: make(map[string]interface{}),
		prefix: config.Map.Prefix,
	}, nil
}

func (ac *mapConnector) validate(_ *Config) error {
	return nil
}
