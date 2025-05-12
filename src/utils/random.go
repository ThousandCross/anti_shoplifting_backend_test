package utils

import (
	rand "crypto/rand"
	"errors"
	mathrand "math/rand"
)

func MakeRandomStr(digit uint32) (string, error) {
	// remove 1,I,l,O, and 0
	const letters = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func MakeRandomFloats(min, max float64) float64 {
	return min + mathrand.Float64()*(max-min)
}
