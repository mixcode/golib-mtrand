package mtrand_test

import (
	"fmt"
	"math/rand"

	"github.com/mixcode/golib/mtrand"
)

func Example_math() {
	// Example of feeding a Mersenne Twister to math/rand
	rng := rand.New(mtrand.NewMT32())
	rng.Seed(1)
	for i := 0; i < 4; i++ {
		fmt.Printf("%08x\n", rng.Int31())
	}

	// Output:
	// 48dbac25
	// 66b09f18
	// 0813e268
	// 0f17f5c4
}
