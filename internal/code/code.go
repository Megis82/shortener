package code

import (
	"crypto/md5"
	"encoding/hex"
)

func CodeString(url string) string {
	hash := md5.Sum([]byte(url))
	shortURL := hex.EncodeToString(hash[:])[:10]
	return shortURL
}
