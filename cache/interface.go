package cache

type ICache interface {
	Get(key, defaultValue string) string
	Set(key, value string) error
}
