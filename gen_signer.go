package taro

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lightninglabs/lndclient"
	"github.com/lightninglabs/taro/asset"
	"github.com/lightningnetwork/lnd/keychain"
)

// LndRpcGenSigner is an implementation of the asset.GenesisSigner interface
// backed by an active lnd node.
type LndRpcGenSigner struct {
	lnd *lndclient.LndServices
}

// NewLndRpcGenSigner returns a new gen signer instance backed by the passed
// connection to a remote lnd node.
func NewLndRpcGenSigner(lnd *lndclient.LndServices) *LndRpcGenSigner {
	return &LndRpcGenSigner{
		lnd: lnd,
	}
}

// SignGenesis signs the passed Genesis description using the public key
// identified by the passed key descriptor. The final tweaked public key and
// the signature are returned.
func (l *LndRpcGenSigner) SignGenesis(keyDesc keychain.KeyDescriptor,
	assetGen asset.Genesis) (*btcec.PublicKey, *schnorr.Signature, error) {

	tweakedPubKey := txscript.ComputeTaprootOutputKey(
		keyDesc.PubKey, assetGen.FamilyKeyTweak(),
	)

	// TODO(roasbeef): needs to actually be a schnorr sig here
	// * fix along side: https://github.com/lightninglabs/taro/issues/33
	// * also make SignMessage suppoprt schnorr
	id := assetGen.ID()
	sig, err := l.lnd.Signer.SignMessage(
		context.Background(), id[:], keyDesc.KeyLocator,
	)
	if err != nil {
		return nil, nil, err
	}

	fakeSchnorrSig, err := schnorr.ParseSignature(sig)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse schnorr sig: "+
			"%w", err)
	}
	return tweakedPubKey, fakeSchnorrSig, nil
}

// A compile time assertion to ensure LndRpcGenSigner meets the
// tarogarden.WalletAnchor interface.
var _ asset.GenesisSigner = (*LndRpcGenSigner)(nil)
