package util

import (
	"math/rand"
	"strings"
)

// RandomInt генерирует случайное целое число в диапазоне [min, max]
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString генерирует случайную строку заданной длины
func RandomString(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner генерирует случайное имя владельца
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney генерирует случайную сумму денег
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomAmount() int64 {
	return RandomInt(-100, 200)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
