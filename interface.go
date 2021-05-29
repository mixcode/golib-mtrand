/*
	interface.go
	interfaces to go standard libraries

	2021-05, github.com/mixcode
*/

package mtrand

// mt32.Seed() is an interface member for math/rand
func (mt *MT32) Seed(seed int64) {
	key := []uint32{uint32(seed & 0xffff_ffff), uint32(seed >> 32)}
	mt.InitByArray(key)
}

// mt32.Int63() is an interface member for math/rand
func (mt *MT32) Int63() int64 {
	r1, r2 := int64(mt.GenUint32()), int64(mt.GenInt31())
	return r2<<32 | r1
}

// mt32.Read() is an io.Reader interface for crypto/random
func (mt *MT32) Read(buf []byte) (n int, err error) {
	l := len(buf)
	for l > 4 {
		u32 := mt.GenUint32()
		buf[0], buf[1], buf[2], buf[3] = byte(u32), byte(u32>>8), byte(u32>>16), byte(u32>>24)
		buf = buf[4:]
		l -= 4
		n += 4
	}
	if l > 0 {
		u32 := mt.GenUint32()
		for i := 0; i < l; i++ {
			buf[i] = byte(u32)
			u32 >>= 8
			n++
		}
	}
	return
}

// mt64.Seed() is an interface member for math/rand
func (mt *MT64) Seed(seed int64) {
	mt.Init(uint64(seed))
}

// mt64.Int63() is an interface member for math/rand
func (mt *MT64) Int63() int64 {
	return mt.GenInt63()
}

// mt64.Read() is an io.Reader interface for crypto/random
func (mt *MT64) Read(buf []byte) (n int, err error) {
	l := len(buf)
	for l > 8 {
		u64 := mt.GenUint64()
		buf[0], buf[1], buf[2], buf[3] = byte(u64), byte(u64>>8), byte(u64>>16), byte(u64>>24)
		buf[4], buf[5], buf[6], buf[7] = byte(u64>>32), byte(u64>>40), byte(u64>>48), byte(u64>>56)
		buf = buf[8:]
		l -= 8
		n += 8
	}
	if l > 0 {
		u64 := mt.GenUint64()
		for i := 0; i < l; i++ {
			buf[i] = byte(u64)
			u64 >>= 8
			n++
		}
	}
	return
}
