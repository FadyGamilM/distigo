package business

type StorageService interface {
	Set(k, v []byte) error
	Get(k []byte) ([]byte, error)
}
