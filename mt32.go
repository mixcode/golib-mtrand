/*

MT32.go is a translation of original Matsumoto-Nishimura implementation of 32-bit Mersenne-Twister Random Number Generator, "mt19937ar.c"

See http://www.math.sci.hiroshima-u.ac.jp/m-mat/MT/MT2002/emt19937ar.html
for original C source code and tech info.

	2021-05, github.com/mixcode


//----------------------------------------------------------------
// Below is the copyright notice of original source code.
//----------------------------------------------------------------

   A C-program for MT19937, with initialization improved 2002/1/26.
   Coded by Takuji Nishimura and Makoto Matsumoto.

   Before using, initialize the state by using init_genrand(seed)
   or init_by_array(init_key, key_length).

   Copyright (C) 1997 - 2002, Makoto Matsumoto and Takuji Nishimura,
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


   Any feedback is very welcome.
   http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt.html
   email: m-mat @ math.sci.hiroshima-u.ac.jp (remove space)
*/

package mtrand

const (
	mt32N         = 624
	mt32M         = 397
	mt32MatrixA   = 0x9908b0df
	mt32UpperMask = 0x8000_0000
	mt32LowerMask = 0x7fff_ffff
)

var (
	mag01 = [2]uint32{0, mt32MatrixA}
)

// 32-bit Mersenne Twister random generator
type MT32 struct {
	mt []uint32 // the array for the state vector
	i  int      // index. if i==mtN+1, then mt[] is not initialized
}

// New() creates a new 32-bit Mersenne Twister random generator
func NewMT32() *MT32 {
	mt := &MT32{mt: make([]uint32, mt32N), i: mt32N + 1}
	return mt
}

// init mt[N] with a seed
func (mt *MT32) Init(seed uint32) {
	mt.mt[0] = seed
	for mt.i = 1; mt.i < mt32N; mt.i++ {
		mt.mt[mt.i] = 1812433253*(mt.mt[mt.i-1]^(mt.mt[mt.i-1]>>30)) + uint32(mt.i)
	}
}

// init with an array
func (mt *MT32) InitByArray(key []uint32) {
	mt.Init(19650218)

	keylen := uint32(len(key))

	k := len(key)
	if mt32N > k {
		k = mt32N
	}

	var i, j uint32 = 1, 0
	for ; k > 0; k-- {
		mt.mt[i] = (mt.mt[i] ^ ((mt.mt[i-1] ^ (mt.mt[i-1] >> 30)) * 1664525)) + key[j] + j // non linear
		i++
		j++
		if i >= mt32N {
			mt.mt[0] = mt.mt[mt32N-1]
			i = 1
		}
		if j >= keylen {
			j = 0
		}
	}
	for k = mt32N - 1; k > 0; k-- {
		mt.mt[i] = (mt.mt[i] ^ ((mt.mt[i-1] ^ (mt.mt[i-1] >> 30)) * 1566083941)) - i // non linear
		i++
		if i >= mt32N {
			mt.mt[0] = mt.mt[mt32N-1]
			i = 1
		}
	}
	mt.mt[0] = 0x8000_0000 // MSB is 1; assuring non-zero initial array
}

// generates a random number on [0,0xffffffff]-interval
func (mt *MT32) GenUint32() uint32 {
	var y uint32

	if mt.i >= mt32N { // generate N words at one time
		var kk int

		if mt.i == mt32N+1 { // if init_genrand() has not been called,
			mt.Init(5489) // a default initial seed is used
		}

		for kk = 0; kk < mt32N-mt32M; kk++ {
			y = (mt.mt[kk] & mt32UpperMask) | (mt.mt[kk+1] & mt32LowerMask)
			mt.mt[kk] = mt.mt[kk+mt32M] ^ (y >> 1) ^ mag01[y&1]
		}
		for ; kk < mt32N-1; kk++ {
			y = (mt.mt[kk] & mt32UpperMask) | (mt.mt[kk+1] & mt32LowerMask)
			mt.mt[kk] = mt.mt[kk+(mt32M-mt32N)] ^ (y >> 1) ^ mag01[y&1]
		}
		y = (mt.mt[mt32N-1] & mt32UpperMask) | (mt.mt[0] & mt32LowerMask)
		mt.mt[mt32N-1] = mt.mt[mt32M-1] ^ (y >> 1) ^ mag01[y&1]

		mt.i = 0
	}

	y = mt.mt[mt.i]
	mt.i++

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	return y
}

// generates a random number on [0,0x7fffffff]-interval
func (mt *MT32) GenInt31() int32 {
	return int32(mt.GenUint32() >> 1)
}

// generates a random number on [0,1]-real-interval
func (mt *MT32) GenReal1() float64 {
	// divided by 2^32-1
	return float64(mt.GenUint32()) * (1.0 / 4294967295.0)
}

// generates a random number on [0,1)-real-interval
func (mt *MT32) GenReal2() float64 {
	// divided by 2^32
	return float64(mt.GenUint32()) * (1.0 / 4294967296.0)
}

// generates a random number on (0,1)-real-interval
func (mt *MT32) GenReal3() float64 {
	// divided by 2^32
	return (float64(mt.GenUint32()) + 0.5) * (1.0 / 4294967296.0)
}

// generates a random number on [0,1) with 53-bit resolution
func (mt *MT32) GenRes53() float64 {
	// These real versions are due to Isaku Wada, 2002/01/09 added
	a, b := mt.GenUint32()>>5, mt.GenUint32()>>6
	return (float64(a)*67108864.0 + float64(b)) * (1.0 / 9007199254740992.0)
}
