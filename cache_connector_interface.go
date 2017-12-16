package cache

type CacheConnectorInterface interface {
	Connect(params map[string]interface{}) StoreInterface

	validate(params map[string]interface{}) map[string]interface{}
}
