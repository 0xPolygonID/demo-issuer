package db

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"go.etcd.io/bbolt"
	"issuer/service/claim"
	"os"
)

var (
	jsonHandle         codec.JsonHandle
	ClaimsBucketName   = []byte("claims")
	IdentityBucketName = []byte("identities")
	ErrKeyNotFound     = fmt.Errorf("key not found")
)

type DB struct {
	conn *bbolt.DB
}

func New(dbFilePath string, removeOldDB bool) (*DB, error) {
	if removeOldDB {
		logger.Info("DB: remove-old-DB flag is true -> delete DB file to start from a clean state")
		_ = os.Remove(dbFilePath)
	}

	logger.Debugf("DB: opening new DB connection (file-path: %s)", dbFilePath)
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
	logger.Trace("DB: init DB")

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

func (db *DB) SaveIdentity(id, authClaimId []byte) error {
	logger.Tracef("DB: saving identity with id: %x", id)

	return db.conn.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(IdentityBucketName).Put(id, authClaimId)
	})
}

func (db *DB) GetClaim(key []byte) (*claim.Claim, error) {
	logger.Tracef("DB: getting claim with the key: %x", key)

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
	logger.Tracef("DB: saving claim with the id: %s", c.ID.String())

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
	logger.Trace("DB: getting all claims")

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
	logger.Trace("DB: getting the saved identity")

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
