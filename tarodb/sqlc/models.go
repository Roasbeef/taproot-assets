// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"database/sql"
	"time"
)

type Addr struct {
	ID               int32
	Version          int16
	GenesisAssetID   int32
	GroupKey         []byte
	ScriptKeyID      int32
	TaprootKeyID     int32
	TaprootOutputKey []byte
	Amount           int64
	AssetType        int16
	CreationTime     time.Time
	ManagedFrom      sql.NullTime
}

type AddrEvent struct {
	ID                  int32
	CreationTime        time.Time
	AddrID              int32
	Status              int16
	ChainTxnID          int32
	ChainTxnOutputIndex int32
	ManagedUtxoID       int32
	AssetProofID        sql.NullInt32
	AssetID             sql.NullInt32
}

type Asset struct {
	AssetID                  int32
	GenesisID                int32
	Version                  int32
	ScriptKeyID              int32
	AssetGroupSigID          sql.NullInt32
	ScriptVersion            int32
	Amount                   int64
	LockTime                 sql.NullInt32
	RelativeLockTime         sql.NullInt32
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	AnchorUtxoID             sql.NullInt32
}

type AssetDelta struct {
	ID                       int32
	OldScriptKey             []byte
	NewAmt                   int64
	NewScriptKey             int32
	SerializedWitnesses      []byte
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	TransferID               int32
	ProofID                  int32
}

type AssetGroup struct {
	GroupID         int32
	TweakedGroupKey []byte
	InternalKeyID   int32
	GenesisPointID  int32
}

type AssetGroupSig struct {
	SigID      int32
	GenesisSig []byte
	GenAssetID int32
	GroupKeyID int32
}

type AssetMintingBatch struct {
	BatchID           int32
	BatchState        int16
	MintingTxPsbt     []byte
	ChangeOutputIndex sql.NullInt32
	GenesisID         sql.NullInt32
	HeightHint        int32
	CreationTimeUnix  time.Time
}

type AssetProof struct {
	ProofID   int32
	AssetID   int32
	ProofFile []byte
}

type AssetSeedling struct {
	SeedlingID      int32
	AssetName       string
	AssetType       int16
	AssetSupply     int64
	AssetMeta       []byte
	EmissionEnabled bool
	BatchID         int32
	GroupGenesisID  sql.NullInt32
}

type AssetTransfer struct {
	ID               int32
	OldAnchorPoint   []byte
	NewInternalKey   int32
	NewAnchorUtxo    int32
	HeightHint       int32
	TransferTimeUnix time.Time
}

type AssetWitness struct {
	WitnessID            int32
	AssetID              int32
	PrevOutPoint         []byte
	PrevAssetID          []byte
	PrevScriptKey        []byte
	WitnessStack         []byte
	SplitCommitmentProof []byte
}

type ChainTxn struct {
	TxnID       int32
	Txid        []byte
	ChainFees   int64
	RawTx       []byte
	BlockHeight sql.NullInt32
	BlockHash   []byte
	TxIndex     sql.NullInt32
}

type GenesisAsset struct {
	GenAssetID     int32
	AssetID        []byte
	AssetTag       string
	MetaData       []byte
	OutputIndex    int32
	AssetType      int16
	GenesisPointID int32
}

type GenesisInfoView struct {
	GenAssetID  int32
	AssetID     []byte
	AssetTag    string
	MetaData    []byte
	OutputIndex int32
	AssetType   int16
	PrevOut     []byte
}

type GenesisPoint struct {
	GenesisID  int32
	PrevOut    []byte
	AnchorTxID sql.NullInt32
}

type InternalKey struct {
	KeyID     int32
	RawKey    []byte
	KeyFamily int32
	KeyIndex  int32
}

type KeyGroupInfoView struct {
	SigID           int32
	GenAssetID      int32
	GenesisSig      []byte
	TweakedGroupKey []byte
	RawKey          []byte
	KeyIndex        int32
	KeyFamily       int32
}

type Macaroon struct {
	ID      []byte
	RootKey []byte
}

type ManagedUtxo struct {
	UtxoID           int32
	Outpoint         []byte
	AmtSats          int64
	InternalKeyID    int32
	TapscriptSibling []byte
	TaroRoot         []byte
	TxnID            int32
}

type MssmtNode struct {
	HashKey   []byte
	LHashKey  []byte
	RHashKey  []byte
	Key       []byte
	Value     []byte
	Sum       int64
	Namespace string
}

type MssmtRoot struct {
	ID        int32
	Namespace string
	RootHash  []byte
}

type ReceiverProofTransferAttempt struct {
	ProofLocatorHash []byte
	TimeUnix         time.Time
}

type ScriptKey struct {
	ScriptKeyID      int32
	InternalKeyID    int32
	TweakedScriptKey []byte
	Tweak            []byte
}

type TransferProof struct {
	ProofID       int32
	TransferID    int32
	SenderProof   []byte
	ReceiverProof []byte
}

type UniverseLeafe struct {
	ID             int32
	AssetGenesisID int32
	UniverseRootID int32
	LeafNodeID     []byte
}

type UniverseRoot struct {
	ID         int32
	RootNodeID int32
	AssetID    []byte
	GroupKey   []byte
}
