package cache

type CacheConnectoInterface interface {
	Connect(params map[string]interface{}) StoreInterface

	validate(params map[string]interface{}) map[string]interface{}
}
