package claims

import (
	"issuer/service/db"
	"issuer/service/identity"
)

type Claims struct {
	identity *identity.Identity
	db       *db.DB
}

func New() (*Claims, error) {
	return &Claims{}, nil
}

func GetClaim() {

}

func SaveClaim() {

}
