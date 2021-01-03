package internal

import (
	"encoding/binary"
)

func PrepareMessage(payload []byte) []byte {
	payloadLen := uint64(len(payload))
	message := make([]byte, payloadLen+ 8) // 8 is the size of uint64

	binary.BigEndian.PutUint64(message, payloadLen)

	if n := copy(message[8:], payload); uint64(n) != payloadLen {
		panic("Error on copying payload data to encode message")
	}

	return message
}

func ExtractPayloadLen(message []byte) uint64 {
	return binary.BigEndian.Uint64(message[0:8])
}