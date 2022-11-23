package identity

import (
	"context"
	"github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"issuer/service/identity/state"
	"issuer/service/models"
	"issuer/utils"
	"math/big"
)

type TransitionInfoResponse struct {
	TxID           string
	BlockTimestamp uint64
	BlockNumber    uint64
}

type TransitionInfoRequest struct {
	IsOldStateGenesis bool
	Identifier        *core.ID
	LatestState       *merkletree.Hash
	NewState          *merkletree.Hash
	Proof             *models.ZKProof
}

type StateStore interface {
	UpdateState(ctx context.Context, trInfo *TransitionInfoRequest) (string, error)
	WaitTransaction(ctx context.Context, txHex string) (*TransitionInfoResponse, error)
}

type Publisher struct {
	i            *Identity
	stateStore   StateStore
	circuitsPath string
}

func (p *Publisher) PrepareInputs() ([]byte, error) {

	// oldState
	oldState, err := circuitsState(p.i.state.CommittedState)
	if err != nil {
		return nil, err
	}

	err = p.i.state.Roots.Tree.Add(context.Background(), oldState.ClaimsRoot.BigInt(), merkletree.HashZero.BigInt())
	if err != nil {
		return nil, err
	}

	newState, err := p.i.state.GetStateHash()
	if err != nil {
		return nil, err
	}

	authInclusionProof, _, err := p.i.state.GetInclusionProof(p.i.authClaim)
	if err != nil {
		return nil, err
	}

	authNonRevocationProof, _, err := p.i.state.GetRevocationProof(p.i.authClaim)
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
		IsOldStateGenesis: p.i.state.CommittedState.IsLatestStateGenesis,

		AuthClaim: authClaim,

		Signature: signature,
	}

	return stateTransitionInputs.InputsMarshal()

}

func (p *Publisher) GenerateProof(ctx context.Context, inputs []byte) (*models.FullProof, error) {

	wasm, err := utils.ReadFileByPath(p.circuitsPath, "/stateTransition/circuit.wasm")
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

	provingKey, err := utils.ReadFileByPath(p.circuitsPath, "/stateTransition/circuit_final.zkey")
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

func (p *Publisher) UpdateState(ctx context.Context, info *TransitionInfoRequest) (string, error) {
	txHex, err := p.stateStore.UpdateState(ctx, info)
	if err != nil {
		return "", err
	}
	go func() {
		tir, err := p.stateStore.WaitTransaction(context.Background(), txHex)
		if err != nil {
			logger.Printf("failed update state from '%s' to '%s'", info.LatestState, info.NewState)
			return
		}
		p.i.state.CommittedState = state.CommittedState{
			Info: &state.Info{
				TxId:           txHex,
				BlockTimestamp: tir.BlockTimestamp,
				BlockNumber:    tir.BlockNumber,
			},

			IsLatestStateGenesis: false,
			RootsTreeRoot:        p.i.state.Roots.Tree.Root(),
			ClaimsTreeRoot:       p.i.state.Claims.Tree.Root(),
			RevocationTreeRoot:   p.i.state.Revocations.Tree.Root(),
		}
	}()
	return txHex, err
}

func circuitsState(s state.CommittedState) (circuits.TreeState, error) {

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
