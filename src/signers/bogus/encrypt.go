package bogus

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/binary"
	"time"

	"tiktok/src/utils"
)

const charset = "Dkdpgh4ZKsQB80/Mfvw36XI1R25-WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe"

func rc4Enc(key, plain []byte) []byte {
	c, _ := rc4.NewCipher(key)
	out := make([]byte, len(plain))
	c.XORKeyStream(out, plain)
	return out
}

func xorReduce(data []byte) byte {
	var x byte
	for _, b := range data {
		x ^= b
	}
	return x
}

func BogusEnc(params string, data string, agent string, cfp int) string {
	agentKey := []byte{0, 1, 8}
	payloadKey := []byte{255}

	paramHash := utils.DoubleMD5([]byte(params))
	dataHash := utils.DoubleMD5([]byte(data))
	agentHash := utils.DoubleMD5([]byte(base64.StdEncoding.EncodeToString(rc4Enc(agentKey, []byte(agent)))))

	buf := make([]byte, 0, 32)
	buf = append(buf, 64)
	buf = append(buf, agentKey...)
	buf = append(buf, paramHash[14:16]...)
	buf = append(buf, dataHash[14:16]...)
	buf = append(buf, agentHash[14:16]...)

	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(time.Now().UnixMilli()))
	buf = append(buf, tmp[:]...)
	binary.BigEndian.PutUint32(tmp[:], uint32(cfp))
	buf = append(buf, tmp[:]...)
	buf = append(buf, xorReduce(buf))

	payload := make([]byte, 0, len(buf)+2)
	payload = append(payload, 2)
	payload = append(payload, payloadKey...)
	payload = append(payload, rc4Enc(payloadKey, buf)...)
	return utils.Base64(payload, charset)
}
