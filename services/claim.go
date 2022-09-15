package services

import (
	"context"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	"lightissuer/models"
	"math/big"
)

const (
	AuthBJJCredentialHash = "ca938857241db9451ea329256b9c06e5"
)

// ClaimService service for handling operations with claims
type ClaimService struct {
	IdentityService *IdentityService
	storage         *DBService
	schemaService   SchemaService
}

// NewClaimService claim service
func NewClaimService(identityService *IdentityService, dbService *DBService, schemaService SchemaService) *ClaimService {
	return &ClaimService{
		identityService,
		dbService,
		schemaService,
	}
}

// Save claim
func (c *ClaimService) Save(claim *models.Claim) (*models.Claim, error) {
	_, err := c.storage.saveClaim(claim)
	if err != nil {
		return nil, errors.Wrap(err, "can't save claim")
	}

	return claim, nil
}

// GetByID get claim for Identity by id
func (c *ClaimService) GetByID(claimID []byte) (*models.Claim, error) {
	var claim *models.Claim
	claim, err := c.storage.getClaim(claimID)
	if err != nil {
		return &models.Claim{}, nil
	}

	return claim, nil
}

// List claims by claims type
func (c *ClaimService) List(identifier *core.ID) ([]*models.Claim, error) {
	var claims []*models.Claim
	var err error

	allClaims, err := c.storage.listClaims()
	if err != nil {
		return nil, err
	}

	for _, claim := range allClaims {
		if claim.Issuer == identifier.String() {
			claims = append(claims, claim)
		}
	}

	return claims, err
}

// GetAuthClaim return Identity authorization claim
func (c *ClaimService) GetAuthClaim(id *core.ID) (*models.Claim, error) {
	claims, err := c.List(id)
	if err != nil {
		return &models.Claim{}, err
	}

	for _, claim := range claims {
		if claim.SchemaHash == AuthBJJCredentialHash {
			return claim, nil
		}
	}

	return &models.Claim{}, errors.New("can't find auth claim for the given ID")
}

// GetRevocationNonceMTP generates MTP proof for given nonce
func (c *ClaimService) GetRevocationNonceMTP(ctx context.Context, nonce uint64) (*verifiable.RevocationStatus, error) {
	rID := new(big.Int).SetUint64(nonce)
	revocationStatus := &verifiable.RevocationStatus{}

	state := c.IdentityService.Identity.State

	revocationStatus.Issuer.State = state.State
	revocationStatus.Issuer.ClaimsTreeRoot = state.ClaimsTreeRoot
	revocationStatus.Issuer.RevocationTreeRoot = state.RevocationTreeRoot
	revocationStatus.Issuer.RootOfRoots = state.RootOfRoots

	if state.RevocationTreeRoot == nil {
		var mtp *merkletree.Proof
		mtp, err := merkletree.NewProofFromData(false, nil, nil)
		if err != nil {
			return nil, err
		}
		revocationStatus.MTP = *mtp
	} else {

		revocationTreeHash, err := merkletree.NewHashFromHex(*state.RevocationTreeRoot)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		identityTrees := c.IdentityService.MTS

		if err != nil {
			return nil, err
		}

		// revocation / non revocation MTP for the latest Identity state
		proof, err := identityTrees.
			GenerateRevocationProof(ctx, rID, revocationTreeHash)
		if err != nil {
			return nil, err
		}

		revocationStatus.MTP = *proof
	}

	return revocationStatus, nil
}
