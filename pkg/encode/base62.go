package encode

import (
	"math"
	"strings"
)

var (
	// base62Str = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	// 为了避免被恶意请求，以将字符串打乱
	base62Str = `x4aU97D0PIWVj3kBdXe2zmorsFGbAtuvcwO5Tq8hiMCfgRY1QJElHKLNnp6SZy`
)

// Int2String 十进制数转为62进制字符串
func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	bl := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, base62Str[mod])
		seq = div
	}

	return string(reverse(bl))
}

// String2Int 62进制字符串转为10进制数
func String2Int(s string) (seq uint64) {
	bl := []byte(s)
	bl = reverse(bl)
	for idx, b := range bl {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return seq
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
