package random

import (
	rand "github.com/qiushiyan/simplebank/foundation/random"
	"github.com/qiushiyan/simplebank/foundation/validate"
)

func RandomOwner() string {
	return rand.RandomString(8)
}

func RandomMoney() int64 {
	return rand.RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{validate.USD, validate.EUR, validate.CAD}
	n := rand.RandomInt(0, int64(len(currencies))-1)
	return currencies[n]
}
