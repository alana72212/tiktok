package utils

func Base64(data []byte, charset string) string {
	var out []byte
	i := 0
	for i+2 < len(data) {
		val := uint32(data[i])<<16 | uint32(data[i+1])<<8 | uint32(data[i+2])
		out = append(out, charset[(val>>18)&63], charset[(val>>12)&63], charset[(val>>6)&63], charset[val&63])
		i += 3
	}
	switch len(data) - i {
	case 1:
		val := uint32(data[i]) << 16
		out = append(out, charset[(val>>18)&63], charset[(val>>12)&63], charset[64], charset[64])
	case 2:
		val := uint32(data[i])<<16 | uint32(data[i+1])<<8
		out = append(out, charset[(val>>18)&63], charset[(val>>12)&63], charset[(val>>6)&63], charset[64])
	}
	return string(out)
}
