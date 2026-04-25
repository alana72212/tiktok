package strdata

import (
	"time"

	"tiktok/src/utils"
)

const charset = "Dkdpgh4ZKsQB80/Mfvw36XI1R25+WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe="

func newPRNG() *utils.PRNG {
	return utils.NewPRNG([4]uint32{uint32(time.Now().UnixMilli()), 0, 0, 0})
}

func nextFloat(p *utils.PRNG) float64 {
	t, r := p.NextWords()
	rv := (r & 0xFFFFF800) >> 11
	return float64(uint64(t)+uint64(rv)) / (1 << 53)
}

func lzw(data []byte) []int {
	table := make(map[string]int)
	for i := 0; i < 256; i++ {
		table[string([]byte{byte(i)})] = i
	}
	next, size := 256, 8
	bucket, filled := 0, 0
	var out []int

	flush := func(code, width int) {
		bucket |= code << filled
		filled += width
		for filled >= 8 {
			out = append(out, bucket&0xFF)
			bucket >>= 8
			filled -= 8
		}
	}

	w := ""
	for _, b := range data {
		wc := w + string([]byte{b})
		if _, exists := table[wc]; exists {
			w = wc
		} else {
			flush(table[w], size)
			table[wc] = next
			next++
			if next > (1 << size) {
				size++
			}
			w = string([]byte{b})
		}
	}
	if w != "" {
		flush(table[w], size)
	}
	if filled > 0 {
		out = append(out, bucket&0xFF)
	}
	return out
}

func StrDataEnc(data string) string {
	p := newPRNG()
	key, kb, rounds := utils.MakeKey(func() float64 { return nextFloat(p) })

	codes := lzw([]byte(data))
	buf := make([]byte, len(codes))
	for i, c := range codes {
		buf[i] = byte(c)
	}

	cipher := utils.Encrypt(key, rounds, buf)
	return utils.Base64(utils.Inter('L', cipher, kb), charset)
}
