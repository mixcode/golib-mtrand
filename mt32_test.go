package mtrand_test

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"

	"testing"

	"github.com/mixcode/golib/mtrand"
)

// compare generated numbers with the output from original program
func TestMT32(t *testing.T) {
	// compare with an output from original C source
	init := []uint32{0x123, 0x234, 0x345, 0x456}
	mt32 := mtrand.NewMT32()
	mt32.InitByArray(init)

	// generate numbers like original test code
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "1000 outputs of genrand_int32()\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(buf, "%10d ", mt32.GenUint32())
		if i%5 == 4 {
			fmt.Fprintf(buf, "\n")
		}
	}
	fmt.Fprintf(buf, "\n1000 outputs of genrand_real2()\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(buf, "%10.8f ", mt32.GenReal2())
		if i%5 == 4 {
			fmt.Fprintf(buf, "\n")
		}
	}
	// split generated output to lines
	generatedLine := make([]string, 0)
	sc := bufio.NewScanner(buf)
	for sc.Scan() {
		generatedLine = append(generatedLine, sc.Text())
	}

	// load output of original C source code
	correctLine := make([]string, 0)
	fi, err := os.Open("original_c_source/mt19937ar/mt19937ar.out")
	if err != nil {
		t.Fatalf("cannot open testdata")
		return
	}
	defer fi.Close()
	// split the text into individual lines to avoid ambiguous linefeed characters
	sc = bufio.NewScanner(fi)
	for sc.Scan() {
		correctLine = append(correctLine, sc.Text())
	}

	// compare lines
	for i := 0; i < len(generatedLine) && i < len(correctLine); i++ {
		if generatedLine[i] != correctLine[i] {
			t.Errorf("value mismatch at line %d", i)
		}
	}
	if len(generatedLine) != len(correctLine) {
		t.Errorf("line count does not match")
	}
}

// test interfaces with go standard lib
func TestMT32Interface(t *testing.T) {

	// first 10 uint32 values with seed 1
	target1 := []uint32{
		0x6ac1f425,
		0xff4780eb,
		0xb8672f8c,
		0xeebc1448,
		0x00077eff,
		0x20ccc389,
		0x4d65aacb,
		0xffc11e85,
		0x2591cb4f,
		0x3c7053c0,
	}
	mt := mtrand.NewMT32()
	mt.Init(1)
	for i, v := range target1 {
		r := mt.GenUint32()
		if r != v {
			t.Errorf("invalid value for iteration %d: expected %x, actual %x", i, v, r)
		}
	}

	// math/rand test
	target2 := []int32{
		0x48dbac25,
		0x66b09f18,
		0x0813e268,
		0x0f17f5c4,
		0x616737a2,
		0x3c728830,
		0x30973b4b,
		0x1adfcc96,
		0x3e721641,
		0x72583673,
	}
	mrand := mrand.New(mtrand.NewMT32())
	mrand.Seed(1)
	for i, v := range target2 {
		r := mrand.Int31()
		if r != v {
			t.Errorf("invalid value for iteration %d: expected %x, actual %x", i, v, r)
		}
	}

	// crypto/rand test
	mt2 := mtrand.NewMT32()
	mt2.Init(1)
	targetPrime := big.NewInt(0)
	targetPrime.SetString("261666422689015313070680887903668429107", 10)
	prime, err := crand.Prime(mt2, 128)
	if err != nil {
		t.Error(err)
	}
	if prime.Cmp(targetPrime) != 0 {
		t.Errorf("generated prime is not equal: explected %v, actual %v", targetPrime, prime)
	}

}
