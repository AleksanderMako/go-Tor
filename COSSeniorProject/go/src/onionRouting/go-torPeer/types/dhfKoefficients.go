package types

import "math/big"

type DFHCoefficients struct {
	N              *big.Int       `json:"n"`
	G              uint64         `json:"g"`
	PublicVariable PublicVariable `json:"publicVariable"`
}
