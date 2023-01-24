package xrand

import (
	c "crypto/rand"
	"encoding/binary"
	"io"
	"math/rand"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number    = "1234567890"
)

var Smol = &SmolRand{}

func init() {
	src := rand.NewSource(0)
	b := make([]byte, 8)
	if _, err := io.ReadFull(c.Reader, b); err != nil {
		panic(err)
	}
	s := binary.BigEndian.Uint64(b)
	src.Seed(int64(s))
	Smol.Rand = rand.New(src)
}

type SmolRand struct {
	Rand *rand.Rand
}

func (sr *SmolRand) GeneratePassword(l int) string {
	res := make([]byte, l)
	idx := 0
	for idx < (l / 3) {
		char := uppercase[sr.Rand.Intn(len(uppercase))]
		res[idx] = char
		idx++
	}
	for idx < (l * 2 / 3) {
		char := number[sr.Rand.Intn(len(number))]
		res[idx] = char
		idx++
	}
	for idx < l {
		char := lowercase[sr.Rand.Intn(len(lowercase))]
		res[idx] = char
		idx++
	}
	sr.Rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return string(res)
}
