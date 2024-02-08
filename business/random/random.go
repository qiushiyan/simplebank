package random

import (
	rand "github.com/qiushiyan/simplebank/foundation/random"
)

func RandomOwner() string {
	return rand.RandomString(8)
}

func RandomMoney() int64 {
	return rand.RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "GBP", "CNY", "JPY", "KRW"}
	n := rand.RandomInt(0, int64(len(currencies))-1)
	return currencies[n]
}
