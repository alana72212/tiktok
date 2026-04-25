package utils

var stateKeys = []uint32{1196819126, 600974999, 3863347763, 1451689750}

func Quarter(s []uint32, a, b, c, d int) {
	s[a] += s[b]
	s[d] ^= s[a]
	s[d] = (s[d] << 16) | (s[d] >> 16)
	s[c] += s[d]
	s[b] ^= s[c]
	s[b] = (s[b] << 12) | (s[b] >> 20)
	s[a] += s[b]
	s[d] ^= s[a]
	s[d] = (s[d] << 8) | (s[d] >> 24)
	s[c] += s[d]
	s[b] ^= s[c]
	s[b] = (s[b] << 7) | (s[b] >> 25)
}

func Block(state []uint32, rounds int) []uint32 {
	x := make([]uint32, 16)
	copy(x, state)
	r := 0
	for r < rounds {
		Quarter(x, 0, 4, 8, 12)
		Quarter(x, 1, 5, 9, 13)
		Quarter(x, 2, 6, 10, 14)
		Quarter(x, 3, 7, 11, 15)
		r++
		if r >= rounds {
			break
		}
		Quarter(x, 0, 5, 10, 15)
		Quarter(x, 1, 6, 11, 12)
		Quarter(x, 2, 7, 12, 13)
		Quarter(x, 3, 4, 13, 14)
		r++
	}
	for i := 0; i < 16; i++ {
		x[i] += state[i]
	}
	return x
}

func Incr(s []uint32) {
	s[12]++
}

func xor(key []uint32, rounds int, data []byte) {
	nFull := len(data) / 4
	words := make([]uint32, (len(data)+3)/4)
	for i := 0; i < nFull; i++ {
		off := 4 * i
		words[i] = uint32(data[off]) | uint32(data[off+1])<<8 | uint32(data[off+2])<<16 | uint32(data[off+3])<<24
	}
	if rem := len(data) % 4; rem != 0 {
		var val uint32
		base := 4 * nFull
		for c := 0; c < rem; c++ {
			val |= uint32(data[base+c]) << (8 * c)
		}
		words[nFull] = val
	}

	state := append([]uint32(nil), key...)
	o := 0
	for o+16 < len(words) {
		ks := Block(state, rounds)
		Incr(state)
		for k := 0; k < 16; k++ {
			words[o+k] ^= ks[k]
		}
		o += 16
	}
	ks := Block(state, rounds)
	for k := 0; k < len(words)-o; k++ {
		words[o+k] ^= ks[k]
	}

	for i := 0; i < nFull; i++ {
		w := words[i]
		off := 4 * i
		data[off] = byte(w)
		data[off+1] = byte(w >> 8)
		data[off+2] = byte(w >> 16)
		data[off+3] = byte(w >> 24)
	}
	if rem := len(data) % 4; rem != 0 {
		w := words[nFull]
		for c := 0; c < rem; c++ {
			data[4*nFull+c] = byte(w >> (8 * c))
		}
	}
}

func Encrypt(key []uint32, rounds int, in []byte) []byte {
	state := make([]uint32, 0, 16)
	state = append(state, stateKeys...)
	state = append(state, key...)
	out := make([]byte, len(in))
	copy(out, in)
	xor(state, rounds, out)
	return out
}
