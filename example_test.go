package mtrand_test

import (
	"fmt"

	"github.com/mixcode/golib/mtrand"
)

func Example() {
	// prepare a 64-bit Mersenne Twister random number generator
	mt64 := mtrand.NewMT64()

	// init with a seed
	mt64.Init(1234567)

	// make some numbers
	fmt.Printf("%v\n", mt64.GenUint64()) // [0, 2^64-1]
	fmt.Printf("%v\n", mt64.GenInt63())  // [0, 2^63-1]
	fmt.Printf("%v\n", mt64.GenReal2())  // [0,1)

	// Output:
	// 18172760479972437302
	// 7074784410711093237
	// 0.656181023402107
}
