package storage

var memStor map[string]string

func Init() {
	memStor = make(map[string]string, 100)
}

func Add(key, value string) {
	memStor[key] = value
}

func Find(key string) string {
	return memStor[key]
}
