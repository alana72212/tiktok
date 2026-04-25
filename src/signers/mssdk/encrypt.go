package mssdk

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func toStr(words []uint32) string {
	out := make([]byte, len(words)*4)
	for i, v := range words {
		off := i * 4
		out[off] = byte(v)
		out[off+1] = byte(v >> 8)
		out[off+2] = byte(v >> 16)
		out[off+3] = byte(v >> 24)
	}
	return string(out)
}

func padKey(a []uint32) []uint32 {
	if len(a) >= 4 {
		return a
	}
	p := make([]uint32, 4)
	copy(p, a)
	return p
}

func toWords(s string, addLen bool) []uint32 {
	l := len(s)
	count := (l + 3) >> 2
	var w []uint32
	if addLen {
		w = make([]uint32, count+1)
		w[count] = uint32(l)
	} else {
		w = make([]uint32, count)
	}
	for i, b := range []byte(s) {
		w[i>>2] |= uint32(b) << ((i & 3) << 3)
	}
	return w
}

func mix(v1, v2, v3 uint32, idx int, e uint32, k []uint32) uint32 {
	a := (v3>>5 ^ v2<<2) + (v2>>3 ^ v3<<4)
	b := (v1 ^ v2) + (k[(idx&3)^int(e)] ^ v3)
	return a ^ b
}

func teaEncrypt(block, key []uint32) []uint32 {
	n := len(block)
	last := n - 1
	acc := block[last]
	var sum uint32
	rounds := 6 + 52/n
	for i := 0; i < rounds; i++ {
		sum += 2654435769
		e := (sum >> 2) & 3
		for j := 0; j < last; j++ {
			acc = block[j] + mix(sum, block[j+1], acc, j, e, key)
			block[j] = acc
		}
		acc = block[last] + mix(sum, block[0], acc, last, e, key)
		block[last] = acc
	}
	return block
}

func encrypt(plain, key string) string {
	cipher := teaEncrypt(toWords(plain, true), padKey(toWords(key, false)))
	res := key + base64.StdEncoding.EncodeToString([]byte(toStr(cipher)))
	res = strings.ReplaceAll(res, "+", "-")
	return strings.ReplaceAll(res, "/", ".")
}

func MssdkEnc(plain string) string {
	key := make([]byte, 4)
	limit := big.NewInt(int64(len(charset)))
	for i := range key {
		n, _ := rand.Int(rand.Reader, limit)
		key[i] = charset[n.Int64()]
	}
	return encrypt(plain, string(key))
}
