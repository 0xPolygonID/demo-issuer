package contract

type GetIdentityResponse struct {
	Identifier string         `codec:"Identifier"`
	State      *IdentityState `codec:"State"`
}

type IdentityState struct {
	StateID            int64  `codec:"StateID"`
	Identifier         string `codec:"Identifier"`
	State              string `codec:"State"`
	RootOfRoots        string `codec:"RootOfRoots"`
	ClaimsTreeRoot     string `codec:"ClaimsTreeRoot"`
	RevocationTreeRoot string `codec:"RevocationTreeRoot"`
}
