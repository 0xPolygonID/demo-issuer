package models

import (
	core "github.com/iden3/go-iden3-core"
)

type Claim struct {
	ID              []byte
	Identifier      *string
	Issuer          string
	SchemaHash      string
	SchemaURL       string
	SchemaType      string
	OtherIdentifier string
	Expiration      int64
	// TODO(illia-korotia): delete from db but left in struct.
	Updatable        bool
	Revoked          bool
	Version          uint32
	RevNonce         uint64
	Data             []byte
	CoreClaim        *core.Claim
	MTPProof         []byte
	SignatureProof   []byte
	IdentityState    *string
	Status           *IdentityStatus
	CredentialStatus []byte
	HIndex           string
}

//// Get returns the value of the core claim
//func (c *CoreClaim) Get() *core.Claim {
//	return (*core.Claim)(c)
//}
//
//// ClaimProof model contains all necessary information to proof claim issuance and later create ZKP base on this info
//type ClaimProof struct {
//	Issuer        string
//	MTPProof      types.JSONText
//	IdentityState IdentityState
//}
//
//// ClaimType internal claim type
//type ClaimType string
//
//// FromClaimer convert Claimer to Claim
//func FromClaimer(claim *core.Claim, schemaURL, schemaType string) (*Claim, error) {
//
//	otherIdentifier := ""
//	id, err := claim.GetID()
//	switch err {
//	case core.ErrNoID:
//	case nil:
//		otherIdentifier = id.String()
//	default:
//		return nil, errors.Wrap(err, "can't get ID")
//	}
//
//	var expiration int64
//	if expirationDate, ok := claim.GetExpirationDate(); ok {
//		expiration = expirationDate.Unix()
//	}
//
//	hindex, err := claim.HIndex()
//	if err != nil {
//		return nil, err
//	}
//
//	sb := claim.GetSchemaHash()
//	schemaHash := hex.EncodeToString(sb[:])
//	res := Claim{
//		SchemaHash:      schemaHash,
//		SchemaURL:       schemaURL,
//		SchemaType:      schemaType,
//		OtherIdentifier: otherIdentifier,
//		Expiration:      expiration,
//		Updatable:       claim.GetFlagUpdatable(),
//		Version:         claim.GetVersion(),
//		RevNonce:        RevNonceUint64(claim.GetRevocationNonce()),
//		CoreClaim:       claim,
//		HIndex:          hindex.String(),
//	}
//
//	return &res, nil
//}
//
//// SetCredentialStatus builds revocation url
//func (c *Claim) SetCredentialStatus(urlBase, issuerID string, sType verifiable.CredentialStatusType) error {
//
//	cStatus := verifiable.CredentialStatus{
//		ID:   fmt.Sprintf("%s/api/v1/claims/revocation/status/%d", urlBase, c.RevNonce),
//		Type: sType,
//	}
//
//	b, err := json.Marshal(cStatus)
//	if err != nil {
//		return err
//	}
//
//	c.CredentialStatus = b
//
//	return nil
//}
//
//// GetCredentialStatus returns CredentialStatus deserialized object
//func (c *Claim) GetCredentialStatus() (*verifiable.CredentialStatus, error) {
//
//	cStatus := &verifiable.CredentialStatus{}
//
//	err := json.Unmarshal(c.CredentialStatus, cStatus)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	return cStatus, nil
//}
//
//// NewCircuitClaimData generates circuits claim structure
//func (c *Claim) NewCircuitClaimData() (circuits.Claim, error) {
//	var circuitClaim circuits.Claim
//	if c == nil {
//		return circuitClaim, errors.New("claim is nil")
//	}
//
//	var err error
//
//	circuitClaim.Claim = c.CoreClaim
//
//	var proof *verifiable.Iden3SparseMerkleProof
//	err = c.MTPProof.Unmarshal(&proof)
//	if err != nil {
//		return circuitClaim, errors.WithStack(err)
//	}
//	circuitClaim.Proof = proof.MTP
//
//	issuerID, err := core.IDFromString(c.Issuer)
//	if err != nil {
//		return circuitClaim, errors.WithStack(err)
//	}
//
//	circuitClaim.IssuerID = &issuerID
//
//	circuitClaim.TreeState = circuits.TreeState{
//		State:          StrMTHex(proof.IssuerData.State.Value),
//		ClaimsRoot:     StrMTHex(proof.IssuerData.State.ClaimsTreeRoot),
//		RevocationRoot: StrMTHex(proof.IssuerData.State.RevocationTreeRoot),
//		RootOfRoots:    StrMTHex(proof.IssuerData.State.RootOfRoots),
//	}
//
//	var sigProof *verifiable.BJJSignatureProof2021
//	err = c.SignatureProof.Unmarshal(&sigProof)
//	if err != nil {
//		return circuitClaim, errors.WithStack(err)
//	}
//
//	signature, err := BJJSignatureFromHexString(sigProof.Signature)
//	if err != nil {
//		return circuitClaim, errors.WithStack(err)
//	}
//
//	circuitClaim.SignatureProof = circuits.BJJSignatureProof{
//		IssuerID: sigProof.IssuerData.ID,
//		IssuerTreeState: circuits.TreeState{
//			State:          StrMTHex(sigProof.IssuerData.State.Value),
//			ClaimsRoot:     StrMTHex(sigProof.IssuerData.State.ClaimsTreeRoot),
//			RevocationRoot: StrMTHex(sigProof.IssuerData.State.RevocationTreeRoot),
//			RootOfRoots:    StrMTHex(sigProof.IssuerData.State.RootOfRoots),
//		},
//		IssuerAuthClaimMTP: sigProof.IssuerData.MTP,
//		Signature:          signature,
//		IssuerAuthClaim:    sigProof.IssuerData.AuthClaim,
//	}
//
//	return circuitClaim, err
//}
//
//// BuildTreeState returns circuits.TreeState structure
//func BuildTreeState(state, claimsTreeRoot, revocationTreeRoot,
//	rootOfRoots *string) (circuits.TreeState, error) {
//
//	return circuits.TreeState{
//		State:          StrMTHex(state),
//		ClaimsRoot:     StrMTHex(claimsTreeRoot),
//		RevocationRoot: StrMTHex(revocationTreeRoot),
//		RootOfRoots:    StrMTHex(rootOfRoots),
//	}, nil
//}
//
//// BJJSignatureFromHexString converts hex to  babyjub.Signature
//func BJJSignatureFromHexString(sigHex string) (*babyjub.Signature, error) {
//	signatureBytes, err := hex.DecodeString(sigHex)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	var sig [64]byte
//	copy(sig[:], signatureBytes)
//	bjjSig, err := new(babyjub.Signature).Decompress(sig)
//	return bjjSig, errors.WithStack(err)
//}
//
//// StrMTHex string to merkle tree hash
//func StrMTHex(s *string) *merkletree.Hash {
//
//	if s == nil {
//		return &merkletree.HashZero
//	}
//
//	h, err := merkletree.NewHashFromHex(*s)
//	if err != nil {
//		return &merkletree.HashZero
//	}
//	return h
//}
