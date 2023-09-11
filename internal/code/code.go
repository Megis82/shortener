package code

import (
	b64 "encoding/base64"
)

func CodeString(input string) string {
	retStr := b64.StdEncoding.EncodeToString([]byte(input))
	return retStr[:10]
}
