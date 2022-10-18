package mtrand_test

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	mtrand "github.com/mixcode/golib-mtrand"
)

func Example_crypto() {
	// Example of feeding Mersenne Twister to crypto/rand.
	// WARNING: please note that the Mersesnne-Twister is NOT cryptographically secure.
	mt64 := mtrand.NewMT64()
	mt64.Init(12345)

	max := big.NewInt(math.MaxInt64)
	n, _ := rand.Int(mt64, max)
	fmt.Println(n)

	// Output:
	// 4182756508441742683
}
