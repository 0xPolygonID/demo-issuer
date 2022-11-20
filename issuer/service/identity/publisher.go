package identity

import (
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"
)

type Publisher struct {
	Identity *Identity
}

func (p *Publisher) PrepareInputs() ([]byte, error) {

	// oldState
	oldState, err := circuitsState(p.Identity.latestRootsState)
	if err != nil {
		return nil, err
	}

	newState, err := p.Identity.state.GetStateHash()
	if err != nil {
		return nil, err
	}

	authInclusionProof, _, err := p.Identity.GetInclusionProof(p.Identity.authClaim)
	if err != nil {
		return nil, err
	}

	authNonRevocationProof, _, err := p.Identity.GetRevocationProof(p.Identity.authClaim)
	if err != nil {
		return nil, err
	}

	authClaim := circuits.Claim{
		Claim:     p.Identity.authClaim,
		TreeState: oldState,
		Proof:     authInclusionProof,
		NonRevProof: &circuits.ClaimNonRevStatus{
			TreeState: oldState,
			Proof:     authNonRevocationProof,
		},
	}

	hashOldAndNewStates, err := poseidon.Hash(
		[]*big.Int{oldState.State.BigInt(), newState.BigInt()})
	if err != nil {
		return nil, err
	}

	signature := p.Identity.sk.SignPoseidon(hashOldAndNewStates)

	stateTransitionInputs := circuits.StateTransitionInputs{
		ID:                p.Identity.Identifier,
		NewState:          newState,
		OldTreeState:      oldState,
		IsOldStateGenesis: p.Identity.latestRootsState.IsLatestStateGenesis,

		AuthClaim: authClaim,

		Signature: signature,
	}

	return stateTransitionInputs.InputsMarshal()

}

func (p *Publisher) GenerateProof(inputs []byte) (*types.ZKProof, error) {

	wasmFile, err := os.Open(wasmPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening wasm file")
	}
	defer wasmFile.Close()

	wasm, err := io.ReadAll(wasmFile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading wasm file")
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

	provingKeyFile, err := os.Open(filepath.Clean(provingKeyPath))
	if err != nil {
		return nil, errors.Wrap(err, "error opening provingKey file")
	}
	defer wasmFile.Close()

	provingKey, err := io.ReadAll(provingKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading provingKey file")
	}

	return prover.Groth16Prover(provingKey, wtnsBytes)
}

func (p *Publisher) SendTx(proof *types.ZKProof) error {
	// TODO://
	return nil

}

func circuitsState(s RootsState) (circuits.TreeState, error) {

	state, err := s.State()
	if err != nil {
		return circuits.TreeState{}, err
	}

	return circuits.TreeState{
		State:          state,
		ClaimsRoot:     s.ClaimsTreeRoot,
		RevocationRoot: s.RevocationTreeRoot,
		RootOfRoots:    s.RootsTreeRoot,
	}, nil
}
