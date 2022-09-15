package services

import (
	"context"
	"github.com/google/uuid"
	client "lightissuer/http"
	"lightissuer/models"
	"lightissuer/utils"

	//nolint:gosec //reason: used for url hash key
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsonSuite "github.com/iden3/go-schema-processor/json"
	jsonldSuite "github.com/iden3/go-schema-processor/json-ld"
	"github.com/iden3/go-schema-processor/loaders"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
)

const (
	Iden3CredentialSchema    = "Iden3Credential"
	Iden3CredentialSchemaURL = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/iden3credential.json-ld"
)

// ErrorLoaderCreatingFailed Error creating loader
var ErrorLoaderCreatingFailed = errors.New("Can't create loader by url")

// CacheService for implementing caching in app
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
}

// SchemaService service for handling operations with claims
type SchemaService struct {
	ipfsURL string
}

// NewSchemaService schema service
func NewSchemaService(ipfsURL string) *SchemaService {
	return &SchemaService{
		ipfsURL: ipfsURL,
	}
}

// Process data and schema and create Index and Value slots
func (s *SchemaService) Process(ctx context.Context, schemaURL, credentialType string, dataBytes []byte) (processor.ParsedSlots, error) {

	loader, err := s.GetLoader(schemaURL)
	if err != nil {
		return processor.ParsedSlots{}, err
	}
	var parser processor.Parser
	var validator processor.Validator
	pr := &processor.Processor{}

	schemaFormat := models.JSONLD
	switch schemaFormat {
	case models.JSON:
		validator = jsonSuite.Validator{}
		parser = jsonSuite.Parser{ParsingStrategy: processor.OneFieldPerSlotStrategy}
	case models.JSONLD:
		validator = jsonldSuite.Validator{ClaimType: credentialType}
		parser = jsonldSuite.Parser{ClaimType: credentialType, ParsingStrategy: processor.OneFieldPerSlotStrategy}
	default:
		return processor.ParsedSlots{}, fmt.Errorf(
			"process suite for %s format is not supported", schemaFormat)
	}

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

// FromClaimToIden3Credential JSON-LD response base on claim
func (s *SchemaService) FromClaimToIden3Credential(claim models.Claim) (*verifiable.Iden3Credential, error) {
	var cred verifiable.Iden3Credential

	cred.ID = claim.ID.String()
	// get typeSpecific credentials schema json

	cred.Context = []string{Iden3CredentialSchemaURL, claim.SchemaURL}

	cred.Type = []string{Iden3CredentialSchema}
	cred.Expiration = time.Unix(claim.Expiration, 0)
	cred.RevNonce = uint64(claim.RevNonce)
	cred.Updatable = claim.Updatable
	cred.Version = claim.Version
	cred.CredentialSchema.ID = claim.SchemaURL
	cred.CredentialSchema.Type = claim.SchemaType

	idp, err := claim.CoreClaim.GetIDPosition()
	if err != nil {
		return nil, err
	}

	cred.SubjectPosition = utils.SubjectPositionIDToString(idp)
	var credData map[string]interface{}

	err = json.Unmarshal([]byte(claim.Data), &credData)
	if err != nil {
		return nil, err
	}
	cred.CredentialSubject = credData
	if claim.OtherIdentifier != "" {
		cred.CredentialSubject["id"] = claim.OtherIdentifier
	}
	cred.CredentialSubject["type"] = claim.SchemaType

	// * create proof object

	proofs := make([]interface{}, 0)

	var signatureProof *verifiable.BJJSignatureProof2021
	if claim.SignatureProof != nil && claim.CredentialStatus.String() != "{}" {
		err = claim.SignatureProof.Unmarshal(&signatureProof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, signatureProof)

	}

	var mtpProof *verifiable.Iden3SparseMerkleProof

	if !strings.EqualFold(claim.MTPProof.String(), "{}") {

		err = claim.MTPProof.Unmarshal(&mtpProof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, mtpProof)

	}
	cred.Proof = proofs

	// create credential status object
	if claim.CredentialStatus != nil && claim.CredentialStatus.String() != "{}" {
		err = claim.CredentialStatus.Unmarshal(&cred.CredentialStatus)
		if err != nil {
			return nil, err
		}
	}
	return &cred, nil
}

// FromIden3CredentialToClaim convert JSON-LD verifiable credentials to Claim model
func (s *SchemaService) FromIden3CredentialToClaim(ctx context.Context, vc verifiable.Iden3Credential) (*models.Claim, error) {

	credentialType := fmt.Sprintf("%v", vc.CredentialSubject["type"])

	parser := jsonldSuite.Parser{
		ClaimType:       credentialType,
		ParsingStrategy: processor.OneFieldPerSlotStrategy,
	}
	schemaURL := vc.CredentialSchema.ID
	schemaType := vc.CredentialSchema.Type
	loader, err := s.GetLoader(schemaURL)
	if err != nil {
		return nil, err
	}
	schemaBytes, _, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}
	coreClaim, err := parser.ParseClaim(&vc, schemaBytes)
	if err != nil {
		return nil, err
	}

	// convert to db model
	claim, err := models.FromClaimer(coreClaim, schemaURL, schemaType)
	if err != nil {
		return nil, err
	}

	data := vc.CredentialSubject
	delete(data, "id")
	delete(data, "type")

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	claim.Data = string(dataJSON)

	claim.ID, err = uuid.Parse(vc.ID)
	if err != nil {
		return nil, err
	}

	claim.Updatable = vc.Updatable
	claim.Expiration = vc.Expiration.Unix()

	claim.CredentialStatus, err = json.Marshal(vc.CredentialStatus)
	if err != nil {
		return nil, err
	}

	// get revocation status
	revStatusBytes, err := client.NewClient(*http.DefaultClient).Get(ctx, vc.CredentialStatus.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs := &models.RevocationStatus{}
	if err = json.Unmarshal(revStatusBytes, rs); err != nil {
		return nil, errors.WithStack(err)
	}
	claim.Revoked = rs.MTP.Existence
	err = parseVerifiableProof(claim, vc.Proof)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

func parseVerifiableProof(claim *models.Claim, vcProof interface{}) error {
	var err error
	switch proofs := vcProof.(type) {
	case []interface{}:
		for _, proof := range proofs {
			err = extractProof(claim, proof)
			if err != nil {
				return err
			}
		}
	case interface{}:
		err = extractProof(claim, vcProof)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractProof(claim *models.Claim, proof interface{}) error {

	switch proof.(map[string]interface{})["@type"] {
	case verifiable.Iden3SparseMerkleProofType:
		sparseMerkleProofBytes, err := json.Marshal(proof)
		if err != nil {
			return err
		}
		var sparseMerkleProof verifiable.Iden3SparseMerkleProof
		err = json.Unmarshal(sparseMerkleProofBytes, &sparseMerkleProof)
		if err != nil {
			return err
		}

		claim.MTPProof = sparseMerkleProofBytes
		claim.Issuer = sparseMerkleProof.IssuerData.ID.String()
		//
	case verifiable.BJJSignatureProofType:
		signatureProofBytes, err := json.Marshal(proof)
		if err != nil {
			return err
		}
		//
		//	// it's an additional check that it can be cast to Signature Proof model
		var signatureProof verifiable.BJJSignatureProof2021
		err = json.Unmarshal(signatureProofBytes, &signatureProof)
		if err != nil {
			return err
		}

		claim.SignatureProof = signatureProofBytes
		claim.Issuer = signatureProof.IssuerData.ID.String()

	default:
		return errors.New("proof type is not supported")
	}
	return nil

}

// Load returns schema content by url
func (s *SchemaService) Load(ctx context.Context, schemaURL string) (schema []byte, extension string, err error) {
	var cacheValue interface{}
	//nolint:gosec //reason: url hash key
	hashBytes := sha1.Sum([]byte(schemaURL))
	hashKey := hex.EncodeToString(hashBytes[:])
	if err != nil {
	}

	// schema doesn't exist in cache. Download and put to cache.
	if cacheValue == nil {
		var loader processor.SchemaLoader
		loader, err = s.GetLoader(schemaURL)
		if err != nil {
			return nil, "", ErrorLoaderCreatingFailed
		}
		var schemaBytes []byte
		schemaBytes, _, err = loader.Load(ctx)
		if err != nil {
			return nil, "", err
		}
		// use request from loader if Redis cache doesn't available.
		return schemaBytes, string(models.JSONLD), nil
	}

	schemaJSONStr, ok := cacheValue.(string)
	if !ok {
		return nil, "", errors.Errorf("can't read schema from cache with url %s and key %s", schemaURL, hashKey)
	}

	return []byte(schemaJSONStr), string(models.JSONLD), nil
}

// GetLoader returns corresponding loader (according to url schema)
func (s *SchemaService) GetLoader(_url string) (processor.SchemaLoader, error) {
	schemaURL, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	switch schemaURL.Scheme {
	case "http", "https":
		return &loaders.HTTP{URL: _url}, nil
	case "ipfs":
		return loaders.IPFS{
			URL: s.ipfsURL,
			CID: schemaURL.Host,
		}, nil
	default:
		return nil, fmt.Errorf("loader for %s is not supported", schemaURL.Scheme)
	}
}
