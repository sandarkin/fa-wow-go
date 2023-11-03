package pow

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

const (
	powHeaderSize = 3
)

type Receiver func(net.Conn) (checkDuration time.Duration, err error)

func NewReceiver(difficulty byte, proofSize int) Receiver {
	bufPool := &sync.Pool{
		New: func() interface{} {
			b := make([]byte, proofSize+powHeaderSize+nonceSize)
			return &b
		},
	}
	return func(conn net.Conn) (checkDuration time.Duration, err error) {
		bufPtr := bufPool.Get().(*[]byte)
		defer bufPool.Put(bufPtr)
		buf := *bufPtr
		resultOffset := powHeaderSize + proofSize
		challengeData := buf[powHeaderSize:resultOffset]

		_, err = rand.Read(challengeData)
		if err != nil {
			return 0, errors.New("failed to read crypto rand")
		}

		buf[0] = difficulty
		binary.BigEndian.PutUint16(buf[1:], uint16(proofSize))

		if _, err = conn.Write(buf[:resultOffset]); err != nil {
			return 0, errors.New("failed to write PoW packet")
		}

		_, err = conn.Read(buf[resultOffset:])
		if err != nil {
			return 0, errors.New("failed to read PoW packet")
		}

		beginCheck := time.Now()
		isValid := CheckBufProof(difficulty, buf[powHeaderSize:])
		checkDuration = time.Since(beginCheck)
		if !isValid {
			return checkDuration, errors.New("is not valid proof")
		}

		return checkDuration, nil
	}
}

func Establish(conn net.Conn) (calcDifficulty byte, calcDuration time.Duration, err error) {
	buf := make([]byte, powHeaderSize)
	_, err = conn.Read(buf)
	if err != nil {
		return 0, 0, errors.New("failed to read PoW header")
	}

	calcDifficulty = buf[0]
	tokenSize := binary.BigEndian.Uint16(buf[1:])

	buf = make([]byte, tokenSize)
	_, err = conn.Read(buf)
	if err != nil {
		err = errors.New("failed to read PoW data")
		return
	}

	beginCalc := time.Now()
	nonce, _, calcErr := CalcProof(calcDifficulty, buf)
	calcDuration = time.Since(beginCalc)
	if calcErr != nil {
		err = calcErr
		return
	}

	_, err = conn.Write(nonce)
	if err != nil {
		err = errors.New("failed to write nonce")
		return
	}

	return
}
