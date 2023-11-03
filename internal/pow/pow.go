package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"math/bits"
)

const nonceSize = 8

func CalcProof(difficulty byte, data []byte) (proofNonce, proofHash []byte, err error) {
	nonceOffset := len(data)
	buf := make([]byte, nonceOffset+nonceSize)
	copy(buf, data)

	var hash [32]byte
	var nonce uint64
	for nonce < math.MaxUint64 {
		binary.BigEndian.PutUint64(buf[nonceOffset:], nonce)
		hash = sha256.Sum256(buf)

		if leadingZerosCount(hash[:]) >= difficulty {
			proofNonce = buf[nonceOffset:]
			proofHash = hash[:]
			return
		} else {
			nonce++
		}
	}
	err = errors.New("unable calculate proof hash")
	return
}

func CheckBufProof(difficulty byte, buf []byte) bool {
	hash := sha256.Sum256(buf)
	return leadingZerosCount(hash[:]) >= difficulty
}

func leadingZerosCount(data []byte) byte {
	count := 0
	for _, v := range data {
		if v == 0 {
			count += 8
		} else {
			count += bits.LeadingZeros8(v)
			break
		}
	}
	if count > 255 {
		return 255
	}
	return byte(count)
}
