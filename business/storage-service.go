package business

type StorageService interface {
	Set(k, v []byte) error
	Get(k []byte) ([]byte, error)
	Pruge(func(key string) bool) error
	LogShardKeys() error
}
