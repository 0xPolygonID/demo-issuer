package utils

import (
	"crypto/rand"
	"encoding/binary"
	core "github.com/iden3/go-iden3-core"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"
)

// SubjectPositionIDToString return string representation for core.IDPosition.
func SubjectPositionIDToString(p core.IDPosition) string {
	switch p {
	case core.IDPositionIndex:
		return SubjectPositionIndex
	case core.IDPositionValue:
		return SubjectPositionValue
	default:
		return ""
	}
}

// RandInt64 generate random uint64
func RandInt64() (uint64, error) {
	var buf [8]byte
	// TODO: this was changed because revocation nonce is cut in dart / js if number is too big
	_, err := rand.Read(buf[:4]) // was rand.Read(buf[:])

	return binary.LittleEndian.Uint64(buf[:]), err
}
