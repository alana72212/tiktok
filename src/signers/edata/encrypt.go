package edata

import (
	"crypto/rand"
	"encoding/base64"

	"tiktok/src/utils"
)

func le32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func chacha20XOR(key [32]byte, nonce [12]byte, src []byte) []byte {
	dst := make([]byte, len(src))

	state := make([]uint32, 16)
	state[0], state[1], state[2], state[3] = 0x61707865, 0x3320646e, 0x79622d32, 0x6b206574
	for i := 0; i < 8; i++ {
		state[4+i] = le32(key[4*i:])
	}
	for i := 0; i < 3; i++ {
		state[13+i] = le32(nonce[4*i:])
	}

	off := 0
	for off < len(src) {
		w := state
		for i := 0; i < 10; i++ {
			utils.Quarter(w[:], 0, 4, 8, 12)
			utils.Quarter(w[:], 1, 5, 9, 13)
			utils.Quarter(w[:], 2, 6, 10, 14)
			utils.Quarter(w[:], 3, 7, 11, 15)
			utils.Quarter(w[:], 0, 5, 10, 15)
			utils.Quarter(w[:], 1, 6, 11, 12)
			utils.Quarter(w[:], 2, 7, 8, 13)
			utils.Quarter(w[:], 3, 4, 9, 14)
		}
		for i := range w {
			w[i] += state[i]
		}

		var ks [64]byte
		for i, v := range w {
			ks[4*i] = byte(v)
			ks[4*i+1] = byte(v >> 8)
			ks[4*i+2] = byte(v >> 16)
			ks[4*i+3] = byte(v >> 24)
		}

		n := len(src) - off
		if n > 64 {
			n = 64
		}
		for i := 0; i < n; i++ {
			dst[off+i] = src[off+i] ^ ks[i]
		}
		off += n
		state[12]++
	}

	return dst
}

func EdataEnc(plaintext string) string {
	var key [32]byte
	var nonce [12]byte
	rand.Read(key[:])
	rand.Read(nonce[:])

	ciphertext := chacha20XOR(key, nonce, []byte(plaintext))

	raw := make([]byte, 0, 1+32+12+len(ciphertext))
	raw = append(raw, 0x01)
	raw = append(raw, key[:]...)
	raw = append(raw, nonce[:]...)
	raw = append(raw, ciphertext...)

	return base64.StdEncoding.EncodeToString(raw)
}
