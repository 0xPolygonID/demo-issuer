package db

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"go.etcd.io/bbolt"
	"issuer/service/claim"
)

var (
	jsonHandle         codec.JsonHandle
	ClaimsBucketName   = []byte("claims")
	IdentityBucketName = []byte("identities")
	ErrKeyNotFound     = fmt.Errorf("key not found")
	IdentityKey        = []byte("identity_key") // TODO: probably needs to be removed
)

type DB struct {
	conn *bbolt.DB
}

func New(dbFilePath string) (*DB, error) {
	//conn, err := bbolt.Open("database.conn", 0600, nil)
	logger.Debugf("opening new DB connection (file-path: %s)", dbFilePath)
	conn, err := bbolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = initDB(conn)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func initDB(conn *bbolt.DB) error {
	return conn.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ClaimsBucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(IdentityBucketName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetConnection() *bbolt.DB {
	return db.conn
}

//func (db *DB) GetIdentity() (string, string, error) {
//	var res string
//
//	return res, db.conn.View(func(tx *bbolt.Tx) error {
//		b := tx.Bucket(IdentityBucketName)
//
//		if val := b.Get(IdentityKey); val != nil {
//			return codec.NewDecoderBytes(val, &jsonHandle).Decode(res)
//		}
//
//		return ErrKeyNotFound
//	})
//}

//func (db *DB) GetIdentity() (*models.Identity, error) {
//	res := &models.Identity{}
//
//	return res, db.conn.View(func(tx *bbolt.Tx) error {
//		b := tx.Bucket(IdentityBucketName)
//
//		if val := b.Get(IdentityKey); val != nil {
//			return codec.NewDecoderBytes(val, &jsonHandle).Decode(res)
//		}
//
//		return ErrKeyNotFound
//	})
//}

func (db *DB) SaveIdentity(id, authClaimId []byte) error {
	return db.conn.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(IdentityBucketName).Put(id, authClaimId)
	})
}

// TODO: remove ID complicated implementation
func (db *DB) GetClaim(key []byte) (*claim.Claim, error) {
	claimB := make([]byte, 0)

	err := db.conn.View(func(tx *bbolt.Tx) error {

		claimB = tx.Bucket(ClaimsBucketName).Get(key)
		if claimB == nil || len(claimB) == 0 {
			return ErrKeyNotFound
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	res := &claim.Claim{}
	err = codec.NewDecoderBytes(claimB, &jsonHandle).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) SaveClaim(c *claim.Claim) error {
	claimB := make([]byte, 0)
	err := codec.NewEncoderBytes(&claimB, &jsonHandle).Encode(c)
	if err != nil {
		return err
	}

	claimIdBytes := []byte(c.ID.String())

	return db.conn.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(ClaimsBucketName).Put(claimIdBytes, claimB)
	})
}

func (db *DB) GetAllClaims() ([]claim.Claim, error) {
	res := []claim.Claim{}

	return res, db.conn.View(func(tx *bbolt.Tx) error {

		b := tx.Bucket(ClaimsBucketName)

		return b.ForEach(func(k, v []byte) error {
			c := claim.Claim{}
			err := codec.NewDecoderBytes(v, &jsonHandle).Decode(c)
			if err != nil {
				return err
			}

			res = append(res, c)
			return nil
		})
	})
}

func (db *DB) GetSavedIdentity() ([]byte, []byte, error) {
	var id []byte
	var authClaimId []byte

	err := db.conn.View(func(tx *bbolt.Tx) error {

		b := tx.Bucket(IdentityBucketName)

		return b.ForEach(func(k, v []byte) error {

			id = k
			authClaimId = v

			return nil
		})
	})
	if err != nil {
		return nil, nil, err
	}

	return id, authClaimId, nil
}
