package gocache

type prefix struct {
	val string
}

func (c *prefix) k(key string) string {
	return c.val + key
}

// Prefix gets the cache key val
func (c *prefix) Prefix() string {
	return c.val
}
