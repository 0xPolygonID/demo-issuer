package db

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"go.etcd.io/bbolt"
	"issuer/service/claim"
	"issuer/service/models"
)

var (
	jsonHandle         codec.JsonHandle
	ClaimsBucketName   = []byte("claims")
	IdentityBucketName = []byte("identities")
	ErrKeyNotFound     = fmt.Errorf("key not found")
	IdentityKey        = []byte("identity") // TODO: probably needs to be removed
)

type DB struct {
	conn *bbolt.DB
}

func New(dbFilePath string) (*DB, error) {
	//conn, err := bbolt.Open("database.conn", 0600, nil)
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

func (db *DB) GetIdentity() (*models.Identity, error) {
	res := &models.Identity{}

	return res, db.conn.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(IdentityBucketName)

		if val := b.Get(IdentityKey); val != nil {
			return codec.NewDecoderBytes(val, &jsonHandle).Decode(res)
		}

		return ErrKeyNotFound
	})
}

func (db *DB) SaveIdentity(iden *models.Identity) error {
	idenB := make([]byte, 0)
	err := codec.NewEncoderBytes(&idenB, &jsonHandle).Encode(iden)
	if err != nil {
		return err
	}

	return db.conn.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(IdentityBucketName).Put(IdentityKey, idenB)
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

	claimId, err := c.CoreClaim.HIndex()
	if err != nil {
		return err
	}

	return db.conn.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(ClaimsBucketName).Put(claimId.Bytes(), claimB)
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
