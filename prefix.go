package gocache

type prefix struct {
	val string
}

func (c *prefix) k(key string) string {
	return c.val + key
}

// GetPrefix gets the cache key val
func (c *prefix) GetPrefix() string {
	return c.val
}
