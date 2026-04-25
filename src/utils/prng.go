package utils

const PRNGRounds = 8

var prngSeed = [12]uint32{
	2517678443, 2718276124, 3212677781, 2633865432,
	217618912, 2931180889, 1498001188, 2157053261,
	211147047, 185100057, 2903579748, 3732962506,
}

type PRNG struct {
	state [16]uint32
	idx   int
}

func NewPRNG(seed [4]uint32) *PRNG {
	p := &PRNG{}
	copy(p.state[:12], prngSeed[:])
	copy(p.state[12:], seed[:])
	return p
}

func (p *PRNG) NextWords() (t, r uint32) {
	bl := Block(p.state[:], PRNGRounds)
	t = bl[p.idx]
	r = bl[p.idx+8]
	if p.idx == 7 {
		Incr(p.state[:])
		p.idx = 0
	} else {
		p.idx++
	}
	return
}
