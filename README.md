
# mtrand: A Mersenne-Twister Random Number Generator

This Go package contains Mersenne-Twister Random Number Generators of 32-bit and 64-bit.

このGoパッケージは32ビット・64ビット版のMersenne-Twister乱数生成器です。

## What is this for?

This package is a translation of original Mersenne-Twister RNGs by Makoto Matsumoto and Takuji Nishimura.
The purpose of this package is to generate exactlay same random number sequences with original RNG implementaion.

このパッケージは広島大の松本先生・西村先生のCバージョンのMersenne Twisterコードを移植したものです。
C関数と同じ乱数列を生成することを目的とします。

## Usage
The package includes two RNG, MT32 and MT64. MT32 is a 32-bit implementaion and MT64 is a 64-bit implementation.

MT32とMT64があります。それぞれ32ビット版・64ビット版です。

#### Example
```
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
```

## Go interface

Additionally, both RNG has interfaces to Go's built-in math/rand and io.Reader for crypto/rand.

Go標準ライブラリのmath/randやcrypto/randに食わせる用のインタフェースも用意されてます。


#### Example of feeding a MT32 to math/rand
```
include "math/rand"
// ...

// use with math/rand
rng := rand.New(mtrand.NewMT32())
rng.Seed(1234)
fmt.Println(rng.Intn(512))
```


## Copyright of original work

See [COPYRIGHT](./COPYRIGHT.md) for copyright notice of original C source codes.

