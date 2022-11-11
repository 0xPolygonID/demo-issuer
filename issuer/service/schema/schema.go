package schema

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	core "github.com/iden3/go-iden3-core"
	jsonldSuite "github.com/iden3/go-schema-processor/json-ld"
	"github.com/iden3/go-schema-processor/loaders"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/pkg/errors"
	"net/url"
)

const (
	Iden3CredentialSchema    = "Iden3Credential"
	Iden3CredentialSchemaURL = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/iden3credential.json-ld"

	AuthBJJCredentialHash = "ca938857241db9451ea329256b9c06e5"
	AuthBJJCredential     = "AuthBJJCredential"
	AuthBJJCredentialURL  = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld"

	// JSONLD JSON-LD schema format
	JSONLD SchemaFormat = "json-ld"

	//// JSON JSON schema format
	//JSON SchemaFormat = "json"

	//IpfsUrl = "ipfs.io"
)

//const (
//	AuthBJJCredentialHash = "ca938857241db9451ea329256b9c06e5"
//	AuthBJJCredential     = "AuthBJJCredential"
//	AuthBJJCredentialURL  = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld"
//)

type SchemaFormat string

type Builder struct {
	ipfsUrl string
}

func NewBuilder(ipfsUrl string) *Builder {
	return &Builder{
		ipfsUrl: ipfsUrl,
	}
}

func (b *Builder) Process(url, _type string, data []byte) (*processor.ParsedSlots, string, error) {
	schemaBytes, _, err := b.load(url)
	if err != nil {
		return nil, "", err
	}

	slots, err := b.getParsedSlots(url, _type, data)
	if err != nil {
		return nil, "", err
	}

	encodedSchema := b.createSchemaHash(schemaBytes, _type)

	return &slots, encodedSchema, nil
}

func (b *Builder) getLoader(_url string) (processor.SchemaLoader, error) {
	schemaURL, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	switch schemaURL.Scheme {
	case "http", "https":
		return &loaders.HTTP{URL: _url}, nil
	case "ipfs":
		return loaders.IPFS{
			URL: b.ipfsUrl,
			CID: schemaURL.Host,
		}, nil
	default:
		return nil, fmt.Errorf("loader for %s is not supported", schemaURL.Scheme)
	}
}

func (b *Builder) getParsedSlots(schemaURL, credentialType string, dataBytes []byte) (processor.ParsedSlots, error) {
	ctx := context.Background()
	loader, err := b.getLoader(schemaURL)
	if err != nil {
		return processor.ParsedSlots{}, err
	}
	var parser processor.Parser
	var validator processor.Validator
	pr := &processor.Processor{}

	// for the case of schemaFormat := "json-ld"
	validator = jsonldSuite.Validator{ClaimType: credentialType}
	parser = jsonldSuite.Parser{ClaimType: credentialType, ParsingStrategy: processor.OneFieldPerSlotStrategy}
	// TODO to remove

	// TODO : it's better to use specific processor (e.g. jsonProcessor.New()), but in this case it's a better option
	pr = processor.InitProcessorOptions(pr, processor.WithValidator(validator), processor.WithParser(parser), processor.WithSchemaLoader(loader))

	schema, _, err := pr.Load(ctx)
	if err != nil {
		return processor.ParsedSlots{}, err
	}
	err = pr.ValidateData(dataBytes, schema)
	if err != nil {
		return processor.ParsedSlots{}, err
	}
	return pr.ParseSlots(dataBytes, schema)
}

func (b *Builder) load(schemaURL string) (schema []byte, extension string, err error) {
	var cacheValue interface{}
	//nolint:gosec //reason: url hash key
	hashBytes := sha1.Sum([]byte(schemaURL))
	hashKey := hex.EncodeToString(hashBytes[:])
	if err != nil {
	}

	// schema doesn't exist in cache. Download and put to cache.
	if cacheValue == nil {
		var loader processor.SchemaLoader
		loader, err = b.getLoader(schemaURL)
		if err != nil {
			return nil, "", err
		}
		var schemaBytes []byte
		schemaBytes, _, err = loader.Load(context.Background())
		if err != nil {
			return nil, "", err
		}
		// use request from loader if Redis cache doesn't available.
		return schemaBytes, string(JSONLD), nil
	}

	schemaJSONStr, ok := cacheValue.(string)
	if !ok {
		return nil, "", errors.Errorf("can't read schema from cache with url %s and key %s", schemaURL, hashKey)
	}

	return []byte(schemaJSONStr), string(JSONLD), nil
}

func (b *Builder) createSchemaHash(schemaBytes []byte, credentialType string) string {
	var sHash core.SchemaHash
	h := crypto.Keccak256(schemaBytes, []byte(credentialType))
	copy(sHash[:], h[len(h)-16:])
	return hex.EncodeToString(sHash[:])
}
