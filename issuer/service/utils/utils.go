package utils

import (
	"encoding/binary"
	"math/rand"
)

// Rand generate random uint64
func Rand() (uint64, error) {
	var buf [8]byte
	// TODO: this was changed because revocation nonce is cut in dart / js if number is too big
	_, err := rand.Read(buf[:4]) // was rand.Read(buf[:])

	return binary.LittleEndian.Uint64(buf[:]), err
}
