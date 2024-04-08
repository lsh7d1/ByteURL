package encode

import (
	"crypto/md5"
	"encoding/hex"
)

// Sum 对传入的参数求md5值
func Sum(data []byte) string {
	return hex.EncodeToString(md5.New().Sum(data)) // 32位16进制数
}
