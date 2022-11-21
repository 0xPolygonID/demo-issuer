// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Proof is an auto generated low-level Go binding around an user-defined struct.
type Proof struct {
	Root     *big.Int
	Siblings [32]*big.Int
	OldKey   *big.Int
	OldValue *big.Int
	IsOld0   bool
	Key      *big.Int
	Value    *big.Int
	Fnc      *big.Int
}

// RootInfo is an auto generated low-level Go binding around an user-defined struct.
type RootInfo struct {
	Root                *big.Int
	ReplacedByRoot      *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtTimestamp *big.Int
	CreatedAtBlock      *big.Int
	ReplacedAtBlock     *big.Int
}

// StateInfo is an auto generated low-level Go binding around an user-defined struct.
type StateInfo struct {
	Id                  *big.Int
	State               *big.Int
	ReplacedByState     *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtTimestamp *big.Int
	CreatedAtBlock      *big.Int
	ReplacedAtBlock     *big.Int
}

// StateMetaData contains all meta data concerning the State contract.
var StateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockN\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getAllStateInfosById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getGISTProof\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_root\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByRoot\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGISTRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_end\",\"type\":\"uint256\"}],\"name\":\"getGISTRootHistory\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGISTRootHistoryLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_root\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfoByBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfoByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getStateInfoById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_state\",\"type\":\"uint256\"}],\"name\":\"getStateInfoByState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"_verifierContractAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newVerifierAddr\",\"type\":\"address\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateEntries\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedBy\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"statesHistories\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_oldState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newState\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_isOldStateGenesis\",\"type\":\"bool\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"transitState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// StateABI is the input ABI used to generate the binding from.
// Deprecated: Use StateMetaData.ABI instead.
var StateABI = StateMetaData.ABI

// State is an auto generated Go binding around an Ethereum contract.
type State struct {
	StateCaller     // Read-only binding to the contract
	StateTransactor // Write-only binding to the contract
	StateFilterer   // Log filterer for contract events
}

// StateCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateSession struct {
	Contract     *State            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateCallerSession struct {
	Contract *StateCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateTransactorSession struct {
	Contract     *StateTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateRaw struct {
	Contract *State // Generic contract binding to access the raw methods on
}

// StateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateCallerRaw struct {
	Contract *StateCaller // Generic read-only contract binding to access the raw methods on
}

// StateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateTransactorRaw struct {
	Contract *StateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewState creates a new instance of State, bound to a specific deployed contract.
func NewState(address common.Address, backend bind.ContractBackend) (*State, error) {
	contract, err := bindState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &State{StateCaller: StateCaller{contract: contract}, StateTransactor: StateTransactor{contract: contract}, StateFilterer: StateFilterer{contract: contract}}, nil
}

// NewStateCaller creates a new read-only instance of State, bound to a specific deployed contract.
func NewStateCaller(address common.Address, caller bind.ContractCaller) (*StateCaller, error) {
	contract, err := bindState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateCaller{contract: contract}, nil
}

// NewStateTransactor creates a new write-only instance of State, bound to a specific deployed contract.
func NewStateTransactor(address common.Address, transactor bind.ContractTransactor) (*StateTransactor, error) {
	contract, err := bindState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateTransactor{contract: contract}, nil
}

// NewStateFilterer creates a new log filterer instance of State, bound to a specific deployed contract.
func NewStateFilterer(address common.Address, filterer bind.ContractFilterer) (*StateFilterer, error) {
	contract, err := bindState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateFilterer{contract: contract}, nil
}

// bindState binds a generic wrapper to an already deployed contract.
func bindState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _State.Contract.StateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _State.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.contract.Transact(opts, method, params...)
}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateCaller) GetAllStateInfosById(opts *bind.CallOpts, _id *big.Int) ([]StateInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getAllStateInfosById", _id)

	if err != nil {
		return *new([]StateInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]StateInfo)).(*[]StateInfo)

	return out0, err

}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateSession) GetAllStateInfosById(_id *big.Int) ([]StateInfo, error) {
	return _State.Contract.GetAllStateInfosById(&_State.CallOpts, _id)
}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateCallerSession) GetAllStateInfosById(_id *big.Int) ([]StateInfo, error) {
	return _State.Contract.GetAllStateInfosById(&_State.CallOpts, _id)
}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTProof(opts *bind.CallOpts, _id *big.Int) (Proof, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTProof", _id)

	if err != nil {
		return *new(Proof), err
	}

	out0 := *abi.ConvertType(out[0], new(Proof)).(*Proof)

	return out0, err

}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateSession) GetGISTProof(_id *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProof(&_State.CallOpts, _id)
}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTProof(_id *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProof(&_State.CallOpts, _id)
}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTProofByBlock(opts *bind.CallOpts, _id *big.Int, _block *big.Int) (Proof, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTProofByBlock", _id, _block)

	if err != nil {
		return *new(Proof), err
	}

	out0 := *abi.ConvertType(out[0], new(Proof)).(*Proof)

	return out0, err

}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateSession) GetGISTProofByBlock(_id *big.Int, _block *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByBlock(&_State.CallOpts, _id, _block)
}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTProofByBlock(_id *big.Int, _block *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByBlock(&_State.CallOpts, _id, _block)
}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTProofByRoot(opts *bind.CallOpts, _id *big.Int, _root *big.Int) (Proof, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTProofByRoot", _id, _root)

	if err != nil {
		return *new(Proof), err
	}

	out0 := *abi.ConvertType(out[0], new(Proof)).(*Proof)

	return out0, err

}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateSession) GetGISTProofByRoot(_id *big.Int, _root *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByRoot(&_State.CallOpts, _id, _root)
}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTProofByRoot(_id *big.Int, _root *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByRoot(&_State.CallOpts, _id, _root)
}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTProofByTime(opts *bind.CallOpts, _id *big.Int, _timestamp *big.Int) (Proof, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTProofByTime", _id, _timestamp)

	if err != nil {
		return *new(Proof), err
	}

	out0 := *abi.ConvertType(out[0], new(Proof)).(*Proof)

	return out0, err

}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateSession) GetGISTProofByTime(_id *big.Int, _timestamp *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByTime(&_State.CallOpts, _id, _timestamp)
}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTProofByTime(_id *big.Int, _timestamp *big.Int) (Proof, error) {
	return _State.Contract.GetGISTProofByTime(&_State.CallOpts, _id, _timestamp)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_State *StateCaller) GetGISTRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_State *StateSession) GetGISTRoot() (*big.Int, error) {
	return _State.Contract.GetGISTRoot(&_State.CallOpts)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_State *StateCallerSession) GetGISTRoot() (*big.Int, error) {
	return _State.Contract.GetGISTRoot(&_State.CallOpts)
}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateCaller) GetGISTRootHistory(opts *bind.CallOpts, _start *big.Int, _end *big.Int) ([]RootInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRootHistory", _start, _end)

	if err != nil {
		return *new([]RootInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]RootInfo)).(*[]RootInfo)

	return out0, err

}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateSession) GetGISTRootHistory(_start *big.Int, _end *big.Int) ([]RootInfo, error) {
	return _State.Contract.GetGISTRootHistory(&_State.CallOpts, _start, _end)
}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_State *StateCallerSession) GetGISTRootHistory(_start *big.Int, _end *big.Int) ([]RootInfo, error) {
	return _State.Contract.GetGISTRootHistory(&_State.CallOpts, _start, _end)
}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_State *StateCaller) GetGISTRootHistoryLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRootHistoryLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_State *StateSession) GetGISTRootHistoryLength() (*big.Int, error) {
	return _State.Contract.GetGISTRootHistoryLength(&_State.CallOpts)
}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_State *StateCallerSession) GetGISTRootHistoryLength() (*big.Int, error) {
	return _State.Contract.GetGISTRootHistoryLength(&_State.CallOpts)
}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTRootInfo(opts *bind.CallOpts, _root *big.Int) (RootInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRootInfo", _root)

	if err != nil {
		return *new(RootInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RootInfo)).(*RootInfo)

	return out0, err

}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateSession) GetGISTRootInfo(_root *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfo(&_State.CallOpts, _root)
}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTRootInfo(_root *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfo(&_State.CallOpts, _root)
}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTRootInfoByBlock(opts *bind.CallOpts, _block *big.Int) (RootInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRootInfoByBlock", _block)

	if err != nil {
		return *new(RootInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RootInfo)).(*RootInfo)

	return out0, err

}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateSession) GetGISTRootInfoByBlock(_block *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfoByBlock(&_State.CallOpts, _block)
}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTRootInfoByBlock(_block *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfoByBlock(&_State.CallOpts, _block)
}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCaller) GetGISTRootInfoByTime(opts *bind.CallOpts, _timestamp *big.Int) (RootInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getGISTRootInfoByTime", _timestamp)

	if err != nil {
		return *new(RootInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RootInfo)).(*RootInfo)

	return out0, err

}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateSession) GetGISTRootInfoByTime(_timestamp *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfoByTime(&_State.CallOpts, _timestamp)
}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCallerSession) GetGISTRootInfoByTime(_timestamp *big.Int) (RootInfo, error) {
	return _State.Contract.GetGISTRootInfoByTime(&_State.CallOpts, _timestamp)
}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCaller) GetStateInfoById(opts *bind.CallOpts, _id *big.Int) (StateInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getStateInfoById", _id)

	if err != nil {
		return *new(StateInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(StateInfo)).(*StateInfo)

	return out0, err

}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateSession) GetStateInfoById(_id *big.Int) (StateInfo, error) {
	return _State.Contract.GetStateInfoById(&_State.CallOpts, _id)
}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCallerSession) GetStateInfoById(_id *big.Int) (StateInfo, error) {
	return _State.Contract.GetStateInfoById(&_State.CallOpts, _id)
}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCaller) GetStateInfoByState(opts *bind.CallOpts, _state *big.Int) (StateInfo, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getStateInfoByState", _state)

	if err != nil {
		return *new(StateInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(StateInfo)).(*StateInfo)

	return out0, err

}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateSession) GetStateInfoByState(_state *big.Int) (StateInfo, error) {
	return _State.Contract.GetStateInfoByState(&_State.CallOpts, _state)
}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_State *StateCallerSession) GetStateInfoByState(_state *big.Int) (StateInfo, error) {
	return _State.Contract.GetStateInfoByState(&_State.CallOpts, _state)
}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_State *StateCaller) GetVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "getVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_State *StateSession) GetVerifier() (common.Address, error) {
	return _State.Contract.GetVerifier(&_State.CallOpts)
}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_State *StateCallerSession) GetVerifier() (common.Address, error) {
	return _State.Contract.GetVerifier(&_State.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateSession) Owner() (common.Address, error) {
	return _State.Contract.Owner(&_State.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_State *StateCallerSession) Owner() (common.Address, error) {
	return _State.Contract.Owner(&_State.CallOpts)
}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_State *StateCaller) StateEntries(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "stateEntries", arg0)

	outstruct := new(struct {
		Id         *big.Int
		Timestamp  *big.Int
		Block      *big.Int
		ReplacedBy *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Block = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.ReplacedBy = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_State *StateSession) StateEntries(arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	return _State.Contract.StateEntries(&_State.CallOpts, arg0)
}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_State *StateCallerSession) StateEntries(arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	return _State.Contract.StateEntries(&_State.CallOpts, arg0)
}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_State *StateCaller) StatesHistories(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "statesHistories", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_State *StateSession) StatesHistories(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _State.Contract.StatesHistories(&_State.CallOpts, arg0, arg1)
}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_State *StateCallerSession) StatesHistories(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _State.Contract.StatesHistories(&_State.CallOpts, arg0, arg1)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateSession) Verifier() (common.Address, error) {
	return _State.Contract.Verifier(&_State.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_State *StateCallerSession) Verifier() (common.Address, error) {
	return _State.Contract.Verifier(&_State.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateTransactor) Initialize(opts *bind.TransactOpts, _verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initialize", _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.Initialize(&_State.TransactOpts, _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_State *StateTransactorSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.Initialize(&_State.TransactOpts, _verifierContractAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateSession) RenounceOwnership() (*types.Transaction, error) {
	return _State.Contract.RenounceOwnership(&_State.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_State *StateTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _State.Contract.RenounceOwnership(&_State.TransactOpts)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_State *StateTransactor) SetVerifier(opts *bind.TransactOpts, _newVerifierAddr common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setVerifier", _newVerifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_State *StateSession) SetVerifier(_newVerifierAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.SetVerifier(&_State.TransactOpts, _newVerifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_State *StateTransactorSession) SetVerifier(_newVerifierAddr common.Address) (*types.Transaction, error) {
	return _State.Contract.SetVerifier(&_State.TransactOpts, _newVerifierAddr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _State.Contract.TransferOwnership(&_State.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_State *StateTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _State.Contract.TransferOwnership(&_State.TransactOpts, newOwner)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactor) TransitState(opts *bind.TransactOpts, _id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "transitState", _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateSession) TransitState(_id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.TransitState(&_State.TransactOpts, _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactorSession) TransitState(_id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.TransitState(&_State.TransactOpts, _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// StateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the State contract.
type StateInitializedIterator struct {
	Event *StateInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StateInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateInitialized represents a Initialized event raised by the State contract.
type StateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) FilterInitialized(opts *bind.FilterOpts) (*StateInitializedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StateInitializedIterator{contract: _State.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StateInitialized) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateInitialized)
				if err := _State.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_State *StateFilterer) ParseInitialized(log types.Log) (*StateInitialized, error) {
	event := new(StateInitialized)
	if err := _State.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the State contract.
type StateOwnershipTransferredIterator struct {
	Event *StateOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StateOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StateOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StateOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateOwnershipTransferred represents a OwnershipTransferred event raised by the State contract.
type StateOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_State *StateFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StateOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _State.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StateOwnershipTransferredIterator{contract: _State.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_State *StateFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StateOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _State.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateOwnershipTransferred)
				if err := _State.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_State *StateFilterer) ParseOwnershipTransferred(log types.Log) (*StateOwnershipTransferred, error) {
	event := new(StateOwnershipTransferred)
	if err := _State.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateStateUpdatedIterator is returned from FilterStateUpdated and is used to iterate over the raw logs and unpacked data for StateUpdated events raised by the State contract.
type StateStateUpdatedIterator struct {
	Event *StateStateUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StateStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStateUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StateStateUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StateStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStateUpdated represents a StateUpdated event raised by the State contract.
type StateStateUpdated struct {
	Id        *big.Int
	BlockN    *big.Int
	Timestamp *big.Int
	State     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateUpdated is a free log retrieval operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_State *StateFilterer) FilterStateUpdated(opts *bind.FilterOpts) (*StateStateUpdatedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return &StateStateUpdatedIterator{contract: _State.contract, event: "StateUpdated", logs: logs, sub: sub}, nil
}

// WatchStateUpdated is a free log subscription operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_State *StateFilterer) WatchStateUpdated(opts *bind.WatchOpts, sink chan<- *StateStateUpdated) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStateUpdated)
				if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStateUpdated is a log parse operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_State *StateFilterer) ParseStateUpdated(log types.Log) (*StateStateUpdated, error) {
	event := new(StateStateUpdated)
	if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
