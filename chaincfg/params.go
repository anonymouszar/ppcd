// Copyright (c) 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"errors"
	"math/big"

	"github.com/ppcsuite/ppcd/wire"
)

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Bitcoin block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)

	// regressionPowLimit is the highest proof of work value a Bitcoin block
	// can have for the regression test network.  It is the value 2^255 - 1.
	regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

	// testNet3PowLimit is the highest proof of work value a Peercoin block
	// can have for the test network (version 3).  It is the value
	// 2^228 - 1.
	testNet3PowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 228), bigOne)

	// simNetPowLimit is the highest proof of work value a Bitcoin block
	// can have for the simulation test network.  It is the value 2^255 - 1.
	simNetPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
)

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
//
// Each checkpoint is selected based upon several factors.  See the
// documentation for btcchain.IsCheckpointCandidate for details on the selection
// criteria.
type Checkpoint struct {
	Height int32
	Hash   *wire.ShaHash
}

// Params defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	Name        string
	Net         wire.BitcoinNet
	DefaultPort string
	DNSSeeds    []string

	// Chain parameters
	GenesisBlock           *wire.MsgBlock
	GenesisMeta            *wire.Meta
	GenesisHash            *wire.ShaHash
	PowLimit               *big.Int
	PowLimitBits           uint32
	SubsidyHalvingInterval int32
	ResetMinDifficulty     bool
	GenerateSupported      bool

	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint

	// Enforce current block version once network has
	// upgraded.  This is part of BIP0034.
	BlockEnforceNumRequired uint64

	// Reject previous block versions once network has
	// upgraded.  This is part of BIP0034.
	BlockRejectNumRequired uint64

	// The number of nodes to check.  This is part of BIP0034.
	BlockUpgradeNumToCheck uint64

	// Mempool parameters
	RelayNonStdTxs bool

	// Address encoding magics
	PubKeyHashAddrID byte // First byte of a P2PKH address
	ScriptHashAddrID byte // First byte of a P2SH address
	PrivateKeyID     byte // First byte of a WIF private key

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32

	// ppc: peercoin specific parameters
	StakeMinAge int64
	// CoinbaseMaturity is the number of blocks required before newly
	// mined bitcoins (coinbase transactions) can be spent.
	CoinbaseMaturity      int64
	InitialHashTargetBits uint32
	// Modifier interval: time to elapse before new modifier is computed
	ModifierInterval         int64
	StakeModifierCheckpoints map[int64]uint32
}

// MainNetParams defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "9901",
	DNSSeeds: []string{
		"seed.peercoin.net",
		"seed2.peercoin.net",
		"seed.peercoin-library.org",
		"ppcseed.ns.7server.net",
	},

	// Chain parameters
	GenesisBlock:           &genesisBlock,
	GenesisHash:            &genesisHash,
	GenesisMeta:            &genesisMeta,
	PowLimit:               mainPowLimit,
	PowLimitBits:           0x1d00ffff,
	SubsidyHalvingInterval: 210000,
	ResetMinDifficulty:     false,
	GenerateSupported:      false,

	// Checkpoints ordered from oldest to newest.
	// https://github.com/peercoin/peercoin/blob/master/src/checkpoints.cpp#L39
	Checkpoints: []Checkpoint{
		{19080, newShaHashFromStr("000000000000bca54d9ac17881f94193fd6a270c1bb21c3bf0b37f588a40dbd7")},
		{30583, newShaHashFromStr("d39d1481a7eecba48932ea5913be58ad3894c7ee6d5a8ba8abeb772c66a6696e")},
		{99999, newShaHashFromStr("27fd5e1de16a4270eb8c68dee2754a64da6312c7c3a0e99a7e6776246be1ee3f")},
		{19999, newShaHashFromStr("ab0dad4b10d2370f009ed6df6effca1ba42f01d5070d6b30afeedf6463fbe7a2")},
		{336000, newShaHashFromStr("4d261cef6e61a5ed8325e560f1d6e36f4698853a4c7134677f47a1d1d842bdf6")},
		{371850, newShaHashFromStr("6b18adcb0a6e080dae85b74eee2b83fabb157bbea64fab0ed2192b2f6d5b89f3")},
		{407813, newShaHashFromStr("00000000000000012730b0f48bed8afbeb08164c9d63597afb082e82ea05cec9")},
		{420000, newShaHashFromStr("fa3fefef369f7f9f0e1b879f42674e8fdfaa88d0172caf1ce67eafed5e684706")},
	},

	// Enforce current block version once majority of the network has
	// upgraded.
	// 75% (750 / 1000)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 95% (950 / 1000)
	BlockEnforceNumRequired: 750,
	BlockRejectNumRequired:  950,
	BlockUpgradeNumToCheck:  1000,

	// Mempool parameters
	RelayNonStdTxs: false,

	// Address encoding magics
	PubKeyHashAddrID: 0x37,        // starts with P
	ScriptHashAddrID: 0x75,        // starts with p
	PrivateKeyID:     0x37 + 0x80, //TODO starts with ? (uncompressed) or ? (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 6,

	// Peercoin
	StakeMinAge:           60 * 60 * 24 * 30, // minimum age for coin age
	CoinbaseMaturity:      100,
	InitialHashTargetBits: 0x1c00ffff,
	ModifierInterval:      6 * 60 * 60, // Set to 6-hour for production network and 20-minute for test network
	StakeModifierCheckpoints: map[int64]uint32{
		0:     uint32(0x0e00670b),
		19080: uint32(0xad4e4d29),
		30583: uint32(0xdc7bf136),
		99999: uint32(0xf555cfd2),
	},
}

// RegressionNetParams defines the network parameters for the regression test
// Bitcoin network.  Not to be confused with the test Bitcoin network (version
// 3), this network is sometimes simply called "testnet".
var RegressionNetParams = Params{
	Name:        "regtest",
	Net:         wire.TestNet,
	DefaultPort: "9903",
	DNSSeeds:    []string{},

	// Chain parameters
	GenesisBlock:           &regTestGenesisBlock,
	GenesisHash:            &regTestGenesisHash,
	GenesisMeta:            &genesisMeta,
	PowLimit:               regressionPowLimit,
	PowLimitBits:           0x207fffff,
	SubsidyHalvingInterval: 150,
	ResetMinDifficulty:     true,
	GenerateSupported:      true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Enforce current block version once majority of the network has
	// upgraded.
	// 75% (750 / 1000)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 95% (950 / 1000)
	BlockEnforceNumRequired: 750,
	BlockRejectNumRequired:  950,
	BlockUpgradeNumToCheck:  1000,

	// Mempool parameters
	RelayNonStdTxs: true,

	// Address encoding magics
	PubKeyHashAddrID: 0x6f, // starts with m or n
	ScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
}

// TestNetParams defines the network parameters for the test Peercoin network
// Not to be confused with the regression test network, this
// network is sometimes simply called "testnet".
var TestNetParams = Params{
	Name:        "testnet",
	Net:         wire.TestNet,
	DefaultPort: "9903",
	DNSSeeds: []string{
		"tseed.peercoin.net",
		"tseed2.peercoin.net",
		"tseed.peercoin-library.org",
	},

	// Chain parameters
	GenesisBlock:           &testNet3GenesisBlock,
	GenesisHash:            &testNet3GenesisHash,
	GenesisMeta:            &genesisMeta,
	PowLimit:               testNet3PowLimit,
	PowLimitBits:           0x1d07ffff,
	SubsidyHalvingInterval: 210000,
	ResetMinDifficulty:     true,
	GenerateSupported:      false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{},

	// Enforce current block version once majority of the network has
	// upgraded.
	// 51% (51 / 100)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 75% (75 / 100)
	BlockEnforceNumRequired: 51,
	BlockRejectNumRequired:  75,
	BlockUpgradeNumToCheck:  100,

	// Mempool parameters
	RelayNonStdTxs: true,

	// Address encoding magics
	PubKeyHashAddrID: 0x6f,        // starts with m or n
	ScriptHashAddrID: 0xc4,        // starts with 2
	PrivateKeyID:     0x6f + 0x80, // TODO(kac-) starts with ? (uncompressed) or ? (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,

	// Peercoin
	StakeMinAge:              60 * 60 * 24, // test net min age is 1 day
	CoinbaseMaturity:         60,
	InitialHashTargetBits:    0x1d07ffff,
	ModifierInterval:         60 * 20, // test net modifier interval is 20 minutes
	StakeModifierCheckpoints: map[int64]uint32{},
}

// SimNetParams defines the network parameters for the simulation test Bitcoin
// network.  This network is similar to the normal test network except it is
// intended for private use within a group of individuals doing simulation
// testing.  The functionality is intended to differ in that the only nodes
// which are specifically specified are used to create the network rather than
// following normal discovery rules.  This is important as otherwise it would
// just turn into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	Net:         wire.SimNet,
	DefaultPort: "18555",
	DNSSeeds:    []string{}, // NOTE: There must NOT be any seeds.

	// Chain parameters
	GenesisBlock:           &simNetGenesisBlock,
	GenesisHash:            &simNetGenesisHash,
	GenesisMeta:            &genesisMeta,
	PowLimit:               simNetPowLimit,
	PowLimitBits:           0x207fffff,
	SubsidyHalvingInterval: 210000,
	ResetMinDifficulty:     true,
	GenerateSupported:      true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Enforce current block version once majority of the network has
	// upgraded.
	// 51% (51 / 100)
	// Reject previous block versions once a majority of the network has
	// upgraded.
	// 75% (75 / 100)
	BlockEnforceNumRequired: 51,
	BlockRejectNumRequired:  75,
	BlockUpgradeNumToCheck:  100,

	// Mempool parameters
	RelayNonStdTxs: true,

	// Address encoding magics
	PubKeyHashAddrID: 0x3f, // starts with S
	ScriptHashAddrID: 0x7b, // starts with s
	PrivateKeyID:     0x64, // starts with 4 (uncompressed) or F (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x00}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3a}, // starts with spub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 115, // ASCII for s
}

var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")

	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)

var (
	registeredNets = map[wire.BitcoinNet]struct{}{
		MainNetParams.Net:       struct{}{},
		TestNetParams.Net:       struct{}{},
		RegressionNetParams.Net: struct{}{},
		SimNetParams.Net:        struct{}{},
	}

	pubKeyHashAddrIDs = map[byte]struct{}{
		MainNetParams.PubKeyHashAddrID: struct{}{},
		TestNetParams.PubKeyHashAddrID: struct{}{}, // shared with regtest
		SimNetParams.PubKeyHashAddrID:  struct{}{},
	}

	scriptHashAddrIDs = map[byte]struct{}{
		MainNetParams.ScriptHashAddrID: struct{}{},
		TestNetParams.ScriptHashAddrID: struct{}{}, // shared with regtest
		SimNetParams.ScriptHashAddrID:  struct{}{},
	}

	// Testnet is shared with regtest.
	hdPrivToPubKeyIDs = map[[4]byte][]byte{
		MainNetParams.HDPrivateKeyID: MainNetParams.HDPublicKeyID[:],
		TestNetParams.HDPrivateKeyID: TestNetParams.HDPublicKeyID[:],
		SimNetParams.HDPrivateKeyID:  SimNetParams.HDPublicKeyID[:],
	}
)

// Register registers the network parameters for a Bitcoin network.  This may
// error with ErrDuplicateNet if the network is already registered (either
// due to a previous Register call, or the network being one of the default
// networks).
//
// Network parameters should be registered into this package by a main package
// as early as possible.  Then, library packages may lookup networks or network
// parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]
	return nil
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsScriptHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsPubKeyHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic
// extended key id and returns the associated public key id.  When the provided
// id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}

	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}

	return pubBytes, nil
}

// newShaHashFromStr converts the passed big-endian hex string into a
// wire.ShaHash.  It only differs from the one available in wire in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newShaHashFromStr(hexStr string) *wire.ShaHash {
	sha, err := wire.NewShaHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
		panic(err)
	}
	return sha
}
