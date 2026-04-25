package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func DoubleMD5(data []byte) [16]byte {
	h := md5.Sum(data)
	return md5.Sum(h[:])
}

func MD5Hex(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}
