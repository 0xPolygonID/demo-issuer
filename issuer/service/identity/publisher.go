package identity

import (
	"math/big"

	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-crypto/poseidon"
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
