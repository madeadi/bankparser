package dictionary

type Storage interface {
	Read(path string) map[string]string
	Write(path string, data map[string]string)
}
