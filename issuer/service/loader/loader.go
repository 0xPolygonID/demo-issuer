package loader

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

const (
	wasmFile         = "circuit.wasm"
	verificationFile = "verification_key.json"
	proofingKeyFile  = "circuit_final.zkey"
)

type Loader struct {
	basePath string
}

func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

func (l *Loader) Wasm(ctx context.Context, circuitName string) ([]byte, error) {
	return uploadFile(l.basePath, circuitName, wasmFile)
}

func (l *Loader) VerificationKey(ctx context.Context, circuitName string) ([]byte, error) {
	return uploadFile(l.basePath, circuitName, verificationFile)
}

func (l *Loader) ProofingKey(ctx context.Context, circuitName string) ([]byte, error) {
	return uploadFile(l.basePath, circuitName, proofingKeyFile)
}

func uploadFile(path, circuitName, fileType string) ([]byte, error) {
	fullPath := filepath.Clean(filepath.Join(path, circuitName, fileType))
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed open file '%s'", fullPath)
	}
	defer func() {
		err := f.Close()
		log.Info("failed close file ", fullPath, err)
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed read file '%s'", fullPath)
	}
	return data, nil
}
