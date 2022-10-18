package mtrand_test

import (
	"bufio"
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"os"

	"testing"

	mtrand "github.com/mixcode/golib-mtrand"
)

// compare generated numbers with the output from original program
func TestMT64(t *testing.T) {
	// compare with an output from original C source
	init := []uint64{0x12345, 0x23456, 0x34567, 0x45678}
	mt64 := mtrand.NewMT64()
	mt64.InitByArray(init)

	// generate numbers like original test code
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "1000 outputs of genrand64_int64()\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(buf, "%20d ", mt64.GenUint64())
		if i%5 == 4 {
			fmt.Fprintf(buf, "\n")
		}
	}
	fmt.Fprintf(buf, "\n1000 outputs of genrand64_real2()\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(buf, "%10.8f ", mt64.GenReal2())
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
	fi, err := os.Open("original_c_source/mt19937-64/mt19937-64.out.txt")
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
			t.Errorf("value mismatch at line %d: %v", i, generatedLine[i])
		}
	}
	if len(generatedLine) != len(correctLine) {
		t.Errorf("line count does not match")
	}
}

// test interfaces with go standard lib
func TestMT64Interface(t *testing.T) {

	// first 10 uint32 values with seed 1
	target1 := []uint64{
		0x2245bd5fbb686f68,
		0x22eb92502318fa4e,
		0x7382d1e77ae6459a,
		0x0561d8057935c08e,
		0x59d47572ecfc6738,
		0xe94ec2d2b9936849,
		0x78833635915bd1b4,
		0x130d84f91bf14b09,
		0x91e180b364f46100,
		0xa29e835c0e448010,
	}
	mt := mtrand.NewMT64()
	mt.Init(1)
	for i, v := range target1 {
		r := mt.GenUint64()
		//_ = i; fmt.Printf("0x%016x,\n", r)
		if r != v {
			t.Errorf("invalid value for iteration %d: expected %x, actual %x", i, v, r)
		}
	}

	// math/rand test
	target2 := []int32{
		0x1122deaf,
		0x1175c928,
		0x39c168f3,
		0x02b0ec02,
		0x2cea3ab9,
		0x74a76169,
		0x3c419b1a,
		0x0986c27c,
		0x48f0c059,
		0x514f41ae,
	}
	mr := mathrand.New(mtrand.NewMT64())
	mr.Seed(1)
	for i, v := range target2 {
		r := mr.Int31()
		//_ = i; fmt.Printf("0x%08x,\n", r)
		if r != v {
			t.Errorf("invalid value for iteration %d: expected %x, actual %x", i, v, r)
		}
	}

	// crypto/rand test
	mt2 := mtrand.NewMT64()
	mt2.Init(1)
	targetN := big.NewInt(0)
	targetN.SetString("7525348656333800738", 10)
	nBig, err := cryptorand.Int(mt2, big.NewInt(math.MaxInt64))

	if err != nil {
		t.Error(err)
	}
	if nBig.Cmp(targetN) != 0 {
		t.Errorf("generated prime is not equal: explected %v, actual %v", targetN, nBig)
	}

}
