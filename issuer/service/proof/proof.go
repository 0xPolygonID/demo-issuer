package proof

import (
	"context"
	"encoding/json"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"
	"issuer/service/loader"
	"issuer/service/models"
)

type Service struct {
	loader *loader.Loader
}

func (p *Service) Generate(ctx context.Context, circuitName string, inputs json.RawMessage) (*models.FullProof, error) {
	wasm, err := p.loader.Wasm(ctx, circuitName)
	if err != nil {
		return nil, err
	}

	calc, err := witness.NewCircom2WitnessCalculator(wasm, true)
	if err != nil {
		return nil, errors.New("can't create witness calculator")
	}

	parsedInputs, err := witness.ParseInputs(inputs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wtnsBytes, err := calc.CalculateWTNSBin(parsedInputs, true)
	if err != nil {
		return nil, errors.New("can't generate witnesses")
	}

	provingKey, err := p.loader.ProofingKey(ctx, circuitName)
	if err != nil {
		return nil, err
	}
	zkpProof, err := prover.Groth16Prover(provingKey, wtnsBytes)
	if err != nil {
		return nil, errors.New("can't generate proof")
	}
	// TODO: get rid of models.Proof structure
	return &models.FullProof{
		Proof: &models.ZKProof{
			A:        zkpProof.Proof.A,
			B:        zkpProof.Proof.B,
			C:        zkpProof.Proof.C,
			Protocol: zkpProof.Proof.Protocol,
		},
		PubSignals: zkpProof.PubSignals,
	}, nil
}
