package mtrand_test

import (
	"crypto/rand"
	"fmt"

	"github.com/mixcode/golib/mtrand"
)

func Example_crypto() {
	// Example of feeding Mersenne Twister to crypto/rand
	mt64 := mtrand.NewMT64()
	mt64.Init(12345)
	prime, _ := rand.Prime(mt64, 128)
	fmt.Println(prime)

	// Output:
	// 334759673973030632160873385234059896671
}
