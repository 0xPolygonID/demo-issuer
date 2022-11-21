package identity

import (
	"context"
	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"
	"issuer/service/blockchain"
	"issuer/service/loader"
	"issuer/service/models"
	"log"
	"math/big"
)

type Publisher struct {
	i          *Identity
	loader     *loader.Loader
	stateStore *blockchain.Blockchain
}

func (p *Publisher) PrepareInputs() ([]byte, error) {

	// oldState
	oldState, err := circuitsState(p.i.latestRootsState)
	if err != nil {
		return nil, err
	}

	newState, err := p.i.state.GetStateHash()
	if err != nil {
		return nil, err
	}

	authInclusionProof, _, err := p.i.GetInclusionProof(p.i.authClaim)
	if err != nil {
		return nil, err
	}

	authNonRevocationProof, _, err := p.i.GetRevocationProof(p.i.authClaim)
	if err != nil {
		return nil, err
	}

	authClaim := circuits.Claim{
		Claim:     p.i.authClaim,
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

	signature := p.i.sk.SignPoseidon(hashOldAndNewStates)

	stateTransitionInputs := circuits.StateTransitionInputs{
		ID:                p.i.Identifier,
		NewState:          newState,
		OldTreeState:      oldState,
		IsOldStateGenesis: p.i.latestRootsState.IsLatestStateGenesis,

		AuthClaim: authClaim,

		Signature: signature,
	}

	return stateTransitionInputs.InputsMarshal()

}

func (p *Publisher) GenerateProof(ctx context.Context, inputs []byte) (*models.FullProof, error) {

	wasm, err := p.loader.Wasm(ctx, "stateTransition")
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

	provingKey, err := p.loader.ProofingKey(ctx, "stateTransition")
	if err != nil {
		return nil, err
	}

	rapidProof, err := prover.Groth16Prover(provingKey, wtnsBytes)
	if err != nil {
		return nil, err
	}

	return &models.FullProof{
		Proof: &models.ZKProof{
			A:        rapidProof.Proof.A,
			B:        rapidProof.Proof.B,
			C:        rapidProof.Proof.C,
			Protocol: rapidProof.Proof.Protocol,
		},
		PubSignals: rapidProof.PubSignals,
	}, nil
}

func (p *Publisher) SendTx(ctx context.Context, info *blockchain.TransitionInfo) (string, error) {
	txHex, err := p.stateStore.UpdateState(ctx, info)
	if err != nil {
		return "", err
	}
	go func() {
		err = p.stateStore.WaitTransaction(context.Background(), txHex)
		if err != nil {
			log.Printf("failed update state from '%s' to '%s'", info.LatestState, info.NewState)
			return
		}
		p.i.latestRootsState = RootsState{
			IsLatestStateGenesis: false,
			RootsTreeRoot:        p.i.state.Roots.Tree.Root(),
			ClaimsTreeRoot:       p.i.state.Claims.Tree.Root(),
			RevocationTreeRoot:   p.i.state.Revocations.Tree.Root(),
		}
	}()
	return txHex, err
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
