package identity

import (
	"github.com/iden3/go-iden3-crypto/babyjub"
	database "issuer/service/db"
	"issuer/service/models"
)

type Identity struct {
	sk       *babyjub.PrivateKey
	pk       *babyjub.PublicKey
	db       *database.DB
	identity *models.Identity
}

func New(db *database.DB, sk *babyjub.PrivateKey, pk *babyjub.PublicKey) (*Identity, error) {
	res := &Identity{
		identity: &models.Identity{},
		db:       db,
		sk:       sk,
		pk:       pk,
	}

	return res, nil
}

func (i *Identity) GetIdentity() string {
	return i.identity.Identifier
}
