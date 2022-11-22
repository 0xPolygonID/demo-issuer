package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	eth "issuer/service/blockchain/contracts"
	"issuer/service/identity"
	"math"
	"math/big"
	"time"
)

type Blockchain struct {
	client          *ethclient.Client
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
}

func NewBlockchainConnect(nodeAddress, contractAddress, pk string) (*Blockchain, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethclient.Dial(nodeAddress)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		client:          ethClient,
		contractAddress: common.HexToAddress(contractAddress),
		privateKey:      privateKey,
	}, nil
}

func (ps *Blockchain) UpdateState(ctx context.Context, trInfo *identity.TransitionInfoRequest) (string, error) {
	if trInfo.NewState.Equals(trInfo.LatestState) {
		return "", errors.New("state hasn't been changed")
	}

	publicKey := ps.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")
	}

	payload, err := ps.getStatePayload(trInfo)
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	tx, err := ps.sendTransaction(ctx, fromAddress, ps.contractAddress, payload)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (ps *Blockchain) WaitTransaction(ctx context.Context, txHex string) (*identity.TransitionInfoResponse, error) {
	txID := common.HexToHash(txHex)
	receipt, err := ps.waitingReceipt(ctx, txID)
	if err != nil {
		return nil, err
	}
	err = ps.waitConfirmation(ctx, txID, receipt.BlockNumber)
	if err != nil {
		return nil, err
	}
	block, err := ps.getBlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, err
	}
	return &identity.TransitionInfoResponse{
		TxID:           txHex,
		BlockTimestamp: block.Time(),
		BlockNumber:    block.NumberU64(),
	}, nil
}

func (ps *Blockchain) waitConfirmation(ctx context.Context, hash common.Hash, formBlock *big.Int) error {
	tryCount := 100
	for tryCount > 0 {
		latestBlock, err := ps.client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		diff := latestBlock - formBlock.Uint64()
		if diff > 10 {
			return nil
		}
		tryCount--
		time.Sleep(time.Second * 5)
	}
	return fmt.Errorf("transaction '%s' is stuck", hash)
}

func (ps *Blockchain) waitingReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	tryCount := 100
	for tryCount > 0 {
		receipt, err := ps.client.TransactionReceipt(ctx, hash)
		if err != nil && errors.Is(err, ethereum.NotFound) {
			logger.Printf("transaction '%s' not found", hash)
			tryCount--
			time.Sleep(time.Second * 5)
			continue
		} else if err != nil {
			return nil, err
		}

		switch receipt.Status {
		case types.ReceiptStatusFailed:
			return nil, fmt.Errorf("transaciton '%s' failed", hash)
		case types.ReceiptStatusSuccessful:
			return receipt, nil
		}
		return nil, fmt.Errorf("unknown tx type '%d'", receipt.Status)
	}
	return nil, fmt.Errorf("all attempts are used")
}

func (ps *Blockchain) getBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return ps.client.BlockByNumber(ctx, number)
}

func (ps *Blockchain) sendTransaction(ctx context.Context, from, to common.Address, payload []byte) (*types.Transaction, error) {
	nonce, err := ps.client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	gasLimit, err := ps.client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from, // the sender of the 'transaction'
		To:    &to,
		Gas:   0,             // wei <-> gas exchange ratio
		Value: big.NewInt(0), // amount of wei sent along with the call
		Data:  payload,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to estimate gas")
	}

	latestBlockHeader, err := ps.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	baseFee := misc.CalcBaseFee(&params.ChainConfig{LondonBlock: big.NewInt(1)}, latestBlockHeader)
	b := math.Round(float64(baseFee.Int64()) * 1.25)
	baseFee = big.NewInt(int64(b))

	gasTip, err := ps.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed get suggest gas tip")
	}

	maxGasPricePerFee := big.NewInt(0).Add(baseFee, gasTip)
	baseTx := &types.DynamicFeeTx{
		To:        &to,
		Nonce:     nonce,
		Gas:       gasLimit,
		Value:     big.NewInt(0),
		Data:      payload,
		GasTipCap: gasTip,
		GasFeeCap: maxGasPricePerFee,
	}

	tx := types.NewTx(baseTx)

	cid, err := ps.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	signer := types.LatestSignerForChainID(cid)

	signedTx, err := types.SignTx(tx, signer, ps.privateKey)
	if err != nil {
		return nil, err
	}

	err = ps.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (ps *Blockchain) getStatePayload(ti *identity.TransitionInfoRequest) ([]byte, error) {
	a, b, c, err := ti.Proof.ProofToBigInts()
	if err != nil {
		return nil, err
	}
	proofA := [2]*big.Int{a[0], a[1]}
	proofB := [2][2]*big.Int{
		{b[0][1], b[0][0]},
		{b[1][1], b[1][0]},
	}
	proofC := [2]*big.Int{c[0], c[1]}

	ab, err := eth.StateMetaData.GetAbi()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	payload, err := ab.Pack(
		"transitState",
		ti.Identifier.BigInt(),
		ti.LatestState.BigInt(),
		ti.NewState.BigInt(),
		ti.IsOldStateGenesis,
		proofA, proofB, proofC)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return payload, nil
}
