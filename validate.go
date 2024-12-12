package signature

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

//go:embed erc6492/contract_sol_ValidateSigOffchain.abi
var validateSigOffchainABI string

//go:embed erc6492/contract_sol_ValidateSigOffchain.bin
var validateSigOffchainBin string

// Bytecode representation of the contract
var universalValidator = common.FromHex(validateSigOffchainBin)

// Contract returns 0x01 if the signature is valid
var valid = common.FromHex("0x01")

// Contract returns 0x00 if the signature is invalid
var invalid = common.FromHex("0x00")

type Validator struct {
	// ContractMetaData is the metadata of the contract
	abi *abi.ABI
	// Client is the client to interact with the contract
	client bind.ContractCaller
}

func NewValidator(client bind.ContractCaller) (*Validator, error) {
	abi, err := abi.JSON(strings.NewReader(validateSigOffchainABI))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}

	return &Validator{
		abi:    &abi,
		client: client,
	}, nil
}

// Validate checks if the signature is valid.
// It returns true if the signature is valid.
// It returns false if the signature is invalid.
func (v *Validator) Validate(ctx context.Context, _signer common.Address, _hash [32]byte, _signature []byte) (bool, error) {
	packed, err := v.abi.Pack("", _signer, _hash, _signature)
	if err != nil {
		return false, fmt.Errorf("abi pack: %w", err)
	}

	buf := bytes.Buffer{}
	buf.Grow(len(universalValidator) + len(packed))
	buf.Write(universalValidator)
	buf.Write(packed)

	result, err := v.client.CallContract(ctx, ethereum.CallMsg{
		Data: buf.Bytes(),
	}, nil)
	if err != nil {
		return false, fmt.Errorf("call contract: %w", err)
	}

	if bytes.Equal(result, valid) {
		return true, nil
	}

	if bytes.Equal(result, invalid) {
		return false, nil
	}

	return false, fmt.Errorf("unexpected result: %x", result)
}
