package xrand

import (
	"crypto/sha512"
	"encoding/binary"
	"math/rand"
)

func Seed(s rand.Source, seed string) {
	n := SumStr(seed)
	s.Seed(int64(n))
}

func SumStr(str string) uint64 {
	b := sha512.Sum512([]byte(str))
	n := binary.BigEndian.Uint64(b[24:])
	n ^= binary.BigEndian.Uint64(b[16:])
	n ^= binary.BigEndian.Uint64(b[8:])
	n ^= binary.BigEndian.Uint64(b[0:])
	return n
}
