package services

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"lightissuer/models"
)

type DBService struct {
	Storage *bbolt.DB
}

var ClaimsBucket = []byte("claims")
var IdentityBucket = []byte("identities")
var ErrKeyNotFound = errors.New("key not found")

// the key used for storing issuer identity { issuer Identity is the only Identity stored in this bucket }
var IdentityKey = []byte("identity")

func NewDBService() (*DBService, error) {
	db, err := newDB()
	if err != nil {
		return &DBService{}, err
	}

	dbService := DBService{
		Storage: db,
	}

	return &dbService, nil
}

func newDB() (*bbolt.DB, error) {

	db, err := bbolt.Open("database.db", 0600, nil)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(ClaimsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(IdentityBucket)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, err
}

// maps public key to Identity!
func (s *DBService) saveIdentity(identity *models.Identity) error {
	err := s.Storage.Update(func(tx *bbolt.Tx) error {
		marshalledIdentity, err := json.Marshal(identity)
		if err != nil {
			return err
		}

		bucket := tx.Bucket(IdentityBucket)

		err = bucket.Put(IdentityKey, marshalledIdentity)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// returns issuer Identity
func (s *DBService) getIdentity() (*models.Identity, error) {

	identity := &models.Identity{}

	err := s.Storage.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket(IdentityBucket)

		val := bucket.Get(IdentityKey)
		if val == nil {
			return ErrKeyNotFound
		}

		err := json.Unmarshal(val, identity)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *DBService) saveClaim(claim *models.Claim) ([]byte, error) {

	key := []byte("")
	claimID := uuid.New()

	err := s.Storage.Update(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket(ClaimsBucket)
		claim.ID = claimID

		marshalledClaimData, err := json.Marshal(claim)
		if err != nil {
			return err
		}

		key = []byte(claimID.String())

		err = bucket.Put(key, marshalledClaimData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *DBService) getClaim(key []byte) (*models.Claim, error) {
	claim := &models.Claim{}

	err := s.Storage.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket(ClaimsBucket)

		val := bucket.Get(key)

		if val == nil {
			return ErrKeyNotFound
		}

		err := json.Unmarshal(val, claim)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &models.Claim{}, err
	}

	claimID, err := uuid.ParseBytes(key)
	if err != nil {
		return &models.Claim{}, err
	}

	claim.ID = claimID

	return claim, nil
}

func (s *DBService) listClaims() ([]*models.Claim, error) {

	var claims []*models.Claim

	err := s.Storage.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ClaimsBucket)

		err := bucket.ForEach(func(k, v []byte) error {
			claim := models.Claim{}
			err := json.Unmarshal(v, &claim)
			if err != nil {
				return err
			}

			claims = append(claims, &claim)

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
