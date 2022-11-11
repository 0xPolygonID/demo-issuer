package models

type Identity struct {
	Identifier         string
	State              string
	RootOfRoots        string
	ClaimsTreeRoot     string
	RevocationTreeRoot string
}
