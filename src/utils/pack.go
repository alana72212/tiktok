package utils

import "encoding/binary"

func MakeKey(next func() float64) (key []uint32, kb []byte, rounds int) {
	key = make([]uint32, 12)
	kb = make([]byte, 48)
	for i := 0; i < 12; i++ {
		w := uint32(next() * 4294967296.0)
		key[i] = w
		rounds = (rounds + int(w&15)) & 15
		binary.LittleEndian.PutUint32(kb[i*4:], w)
	}
	rounds += 5
	return
}

func Inter(prefix byte, cipher, kb []byte) []byte {
	n := len(cipher) + 1
	pos := 0
	for _, b := range kb {
		pos = (pos + int(b)) % n
	}
	for _, b := range cipher {
		pos = (pos + int(b)) % n
	}
	out := make([]byte, 0, 1+len(cipher)+len(kb))
	out = append(out, prefix)
	out = append(out, cipher[:pos]...)
	out = append(out, kb...)
	out = append(out, cipher[pos:]...)
	return out
}
