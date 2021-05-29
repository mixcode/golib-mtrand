/*

Package mtrand is Mersenne-Twister random number generators of 32-bit and 64-bit.

This package is a translation of original Mersenne-Twister RNGs by Makoto Matsumoto and Takuji Nishimura, of Hiroshima Univ.

The purpose of this package is to provide RNGs that generate exact same random number sequences with the original implementaion. Functions resembles the functions in original, but names slightly changed for Go convension.
(See following link for the original works http://www.math.sci.hiroshima-u.ac.jp/m-mat/MT/emt.html)

MT32 is 32-bit Mersenne Twister RNG, which generates same sequences with "mt19937ar.c" implementation.
MT64 is 64-bit Mersenne Twister RNG, which generates same sequences with "mt19937-64.c" implementation as well.

Additionally, both RNGs have interfaces for Go's built-in math/rand and cryto/rand.


*/
package mtrand
