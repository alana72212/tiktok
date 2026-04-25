package gnarly

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"time"

	"tiktok/src/utils"
)

const charset = "u09tbS3UvgDEe6r-ZVMXzLpsAohTn7mdINQlW412GqBjfYiyk8JORCF5/xKHwacP="

func randU32() uint32 {
	v, _ := rand.Int(rand.Reader, big.NewInt(1<<32))
	return uint32(v.Uint64())
}

func newPRNG() *utils.PRNG {
	return utils.NewPRNG([4]uint32{
		uint32(time.Now().UnixMilli()),
		randU32(), randU32(), randU32(),
	})
}

func nextFloat(p *utils.PRNG) float64 {
	t, r := p.NextWords()
	rv := (r & 0xFFFFFFF0) >> 11
	return (float64(t) + 4294967296.0*float64(rv)) / (1 << 53)
}

func toBytes(v uint32) []byte {
	if v < 65025 {
		out := [2]byte{}
		binary.BigEndian.PutUint16(out[:], uint16(v))
		return out[:]
	}
	out := [4]byte{}
	binary.BigEndian.PutUint32(out[:], v)
	return out[:]
}

func GnarlyEnc(query string, body string, agent string, ver string, cfp int) string {
	ts := uint32(time.Now().Unix())
	ts2 := uint32((time.Now().UnixMilli()%2147483648 + 1) & 0x7FFFFFFF)
	obj := map[int]interface{}{
		1: uint32(1), 2: 0,
		3: utils.MD5Hex([]byte(query)),
		4: utils.MD5Hex([]byte(body)),
		5: utils.MD5Hex([]byte(agent)),
		6: ts, 7: uint32(cfp), 8: ts2, 9: ver,
	}

	var order []int
	switch ver {
	case "5.1.0":
		order = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	case "5.1.1":
		obj[10] = "2.0.0.430"
		obj[11] = uint32(1)
		var v12 uint32
		for i := 1; i <= 11; i++ {
			switch v := obj[i].(type) {
			case uint32:
				v12 ^= v
			case string:
				var tmp uint32
				for _, b := range []byte(v)[:min(4, len(v))] {
					tmp = (tmp << 8) | uint32(b)
				}
				v12 ^= tmp
			}
		}
		obj[12] = v12
		order = []int{7, 1, 11, 5, 12, 6, 10, 9, 2, 4, 3, 8, 0}
	default:
		panic("unsupported version: " + ver)
	}

	var v0 uint32
	last := len(obj)
	for i := 1; i <= last; i++ {
		if v, ok := obj[i].(uint32); ok {
			v0 ^= v
		}
	}
	obj[0] = v0

	pay := []byte{byte(len(obj))}
	for _, idx := range order {
		pay = append(pay, byte(idx))
		var vb []byte
		switch t := obj[idx].(type) {
		case uint32:
			vb = toBytes(t)
		case string:
			vb = []byte(t)
		}
		pay = append(pay, toBytes(uint32(len(vb)))...)
		pay = append(pay, vb...)
	}

	p := newPRNG()
	key, kb, rounds := utils.MakeKey(func() float64 { return nextFloat(p) })
	cipher := utils.Encrypt(key, rounds, pay)
	return utils.Base64(utils.Inter('K', cipher, kb), charset)
}
