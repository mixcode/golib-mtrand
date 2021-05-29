/*

mt64.go is a translation of original Matsumoto-Nishimura implementation of 64-bit Mersenne-Twister Random Number Generator, "mt19937-64.c"

See http://www.math.sci.hiroshima-u.ac.jp/m-mat/MT/emt64.html
for original C source code and tech info.

	2021-05, github.com/mixcode


//----------------------------------------------------------------
// Below is the copyright notice of original source code.
//----------------------------------------------------------------

   A C-program for MT19937-64 (2004/9/29 version).
   Coded by Takuji Nishimura and Makoto Matsumoto.

   This is a 64-bit version of Mersenne Twister pseudorandom number
   generator.

   Before using, initialize the state by using init_genrand64(seed)
   or init_by_array64(init_key, key_length).

   Copyright (C) 2004, Makoto Matsumoto and Takuji Nishimura,
   All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:

     1. Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.

     2. Redistributions in binary form must reproduce the above copyright
        notice, this list of conditions and the following disclaimer in the
        documentation and/or other materials provided with the distribution.

     3. The names of its contributors may not be used to endorse or promote
        products derived from this software without specific prior written
        permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR
   CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
   EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
   PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
   PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
   LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
   SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

   References:
   T. Nishimura, ``Tables of 64-bit Mersenne Twisters''
     ACM Transactions on Modeling and
     Computer Simulation 10. (2000) 348--357.
   M. Matsumoto and T. Nishimura,
     ``Mersenne Twister: a 623-dimensionally equidistributed
       uniform pseudorandom number generator''
     ACM Transactions on Modeling and
     Computer Simulation 8. (Jan. 1998) 3--30.

   Any feedback is very welcome.
   http://www.math.hiroshima-u.ac.jp/~m-mat/MT/emt.html
   email: m-mat @ math.sci.hiroshima-u.ac.jp (remove spaces)
*/

package mtrand

const (
	mt64NN      = 312
	mt64MM      = 156
	mt64MatrixA = 0xB502_6F5A_A966_19E9
	mt64UM      = 0xFFFF_FFFF_8000_0000 // Most significant 33 bits
	mt64LM      = 0x0000_0000_7FFF_FFFF // Least significant 31 bits
)

var (
	mt64mag01 = []uint64{0, mt64MatrixA}
)

// 64-bit Mersenne Twister random generator
type MT64 struct {
	mt []uint64 // the array for the state vector
	i  int      // index. if i==mt64NN+1, then mt[] is not initialized
}

// New() creates a new 32-bit Mersenne Twister random generator
func NewMT64() *MT64 {
	mt := &MT64{mt: make([]uint64, mt64NN), i: mt64NN + 1}
	return mt
}

// initializes mt[mt64NN] with a seed
func (mt *MT64) Init(seed uint64) {
	mt.mt[0] = seed
	for mt.i = 1; mt.i < mt64NN; mt.i++ {
		mt.mt[mt.i] = 6364136223846793005*(mt.mt[mt.i-1]^(mt.mt[mt.i-1]>>62)) + uint64(mt.i)
	}
}

// initialize by an array with array-length
// init_key is the array for initializing keys
// key_length is its length
func (mt *MT64) InitByArray(init_key []uint64) {
	mt.Init(19650218)

	key_length := uint64(len(init_key))

	k := len(init_key)
	if mt64NN > k {
		k = mt64NN
	}

	var i, j uint64 = 1, 0
	for ; k > 0; k-- {
		mt.mt[i] = (mt.mt[i] ^ ((mt.mt[i-1] ^ (mt.mt[i-1] >> 62)) * 3935559000370003845)) + init_key[j] + j // non linear
		i++
		j++
		if i >= mt64NN {
			mt.mt[0] = mt.mt[mt64NN-1]
			i = 1
		}
		if j >= key_length {
			j = 0
		}
	}
	for k = mt64NN - 1; k > 0; k-- {
		mt.mt[i] = mt.mt[i] ^ ((mt.mt[i-1] ^ (mt.mt[i-1] >> 62)) * 2862933555777941757) - i // non linear
		i++
		if i >= mt64NN {
			mt.mt[0] = mt.mt[mt64NN-1]
			i = 1
		}
	}

	mt.mt[0] = 1 << 63 // MSB is 1; assuring non-zero initial array
}

// generates a random number on [0, 2^64-1]-interval
func (mt *MT64) GenUint64() uint64 {
	var i int
	var x uint64

	if mt.i >= mt64NN { // generate mt64NN words at one time

		// if init_genrand64() has not been called,
		// a default initial seed is used
		if mt.i == mt64NN+1 {
			mt.Init(5489)
		}

		for i = 0; i < mt64NN-mt64MM; i++ {
			x = (mt.mt[i] & mt64UM) | (mt.mt[i+1] & mt64LM)
			mt.mt[i] = mt.mt[i+mt64MM] ^ (x >> 1) ^ mt64mag01[x&1]
		}
		for ; i < mt64NN-1; i++ {
			x = (mt.mt[i] & mt64UM) | (mt.mt[i+1] & mt64LM)
			mt.mt[i] = mt.mt[i+(mt64MM-mt64NN)] ^ (x >> 1) ^ mt64mag01[x&1]
		}
		x = (mt.mt[mt64NN-1] & mt64UM) | (mt.mt[0] & mt64LM)
		mt.mt[mt64NN-1] = mt.mt[mt64MM-1] ^ (x >> 1) ^ mt64mag01[x&1]

		mt.i = 0
	}

	x = mt.mt[mt.i]
	mt.i++

	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)

	return x
}

// generates a random number on [0, 2^63-1]-interval
func (mt *MT64) GenInt63() int64 {
	return (int64)(mt.GenUint64() >> 1)
}

// generates a random number on [0,1]-real-interval
func (mt *MT64) GenReal1() float64 {
	return float64(mt.GenUint64()>>11) * (1.0 / 9007199254740991.0)
}

// generates a random number on [0,1)-real-interval
func (mt *MT64) GenReal2() float64 {
	return float64(mt.GenUint64()>>11) * (1.0 / 9007199254740992.0)
}

// generates a random number on (0,1)-real-interval
func (mt *MT64) GenReal3() float64 {
	return (float64(mt.GenUint64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}
