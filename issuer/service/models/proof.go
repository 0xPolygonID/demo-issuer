package models

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
)

// ZKProof is structure that represents SnarkJS library result of proof generation
type ZKProof struct {
	A        []string   `json:"pi_a"`
	B        [][]string `json:"pi_b"`
	C        []string   `json:"pi_c"`
	Protocol string     `json:"protocol"`
}

// ProofToBigInts transforms a zkp (that uses `*bn256.G1` and `*bn256.G2`) into
// `*big.Int` format, to be used for example in snarkjs solidity verifiers.
func (p *ZKProof) ProofToBigInts() (a []*big.Int, b [][]*big.Int, c []*big.Int, err error) {

	a, err = ArrayStringToBigInt(p.A)
	if err != nil {
		return nil, nil, nil, err
	}
	b = make([][]*big.Int, len(p.B))
	for i, v := range p.B {
		b[i], err = ArrayStringToBigInt(v)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	c, err = ArrayStringToBigInt(p.C)
	if err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

// FullProof is ZKP proof with public signals
type FullProof struct {
	Proof      *ZKProof `json:"proof"`
	PubSignals []string `json:"pub_signals"`
}

// ArrayStringToBigInt converts array of string to big int
func ArrayStringToBigInt(s []string) ([]*big.Int, error) {
	var o []*big.Int
	for i := 0; i < len(s); i++ {
		si, err := stringToBigInt(s[i])
		if err != nil {
			return o, nil
		}
		o = append(o, si)
	}
	return o, nil
}

func stringToBigInt(s string) (*big.Int, error) {
	base := 10
	if bytes.HasPrefix([]byte(s), []byte("0x")) {
		base = 16
		s = strings.TrimPrefix(s, "0x")
	}
	n, ok := new(big.Int).SetString(s, base)
	if !ok {
		return nil, fmt.Errorf("can not parse string to *big.Int: %s", s)
	}
	return n, nil
}
