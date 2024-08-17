package rfq

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
)

// arithScale is the scale used for arithmetic operations. This is used to
// ensure that we don't lose precision when doing arithmetic operations.
const arithScale = 11

// MilliSatoshiToUnits converts the given milli-satoshi amount to units using
// the given price in units per bitcoin as a fixed point in the asset's desired
// resolution (scale equal to decimal display).
//
// Given the amount of mSat (X), and the number of units per BTC (Y), we can
// compute the total amount of units (U) as follows:
//   - U = (X / M) * Y
//   - where M is the number of mSAT in a BTC (100,000,000,000).
func MilliSatoshiToUnits[N Int[N]](milliSat lnwire.MilliSatoshi,
	unitsPerBtc FixedPoint[N]) FixedPoint[N] {

	// TODO(roasbeef): take max of arith scale and scale of units per btc

	// Before we do any computation, we'll scale everything up to our
	// arithmetric scale.
	mSatFixed := FixedPointFromUint64[N](
		uint64(milliSat), arithScale,
	)
	scaledUnitsPerBtc := unitsPerBtc.ScaleTo(arithScale)

	// Next, we'll convert the amount of mSAT to BTC. We do this by
	// dividing by the number of mSAT in a BTC.
	oneBtcInMilliSat := FixedPointFromUint64[N](
		uint64(btcutil.SatoshiPerBitcoin*1_000), arithScale,
	)
	amtBTC := mSatFixed.Div(oneBtcInMilliSat)

	// Now that we have the amount of BTC as input, and the amount of units
	// per BTC, we multiply the two to get the total amount of units.
	amtUnits := amtBTC.Mul(scaledUnitsPerBtc)

	// The final response will need to scale back down to the original
	// amount of units that were passed in.
	scaledAmt := amtUnits.ScaleTo(unitsPerBtc.Scale)

	return scaledAmt
}

// UnitsToMilliSatoshi converts the given number of asset units to a
// milli-satoshi amount, using the given price in units per bitcoin as a fixed
// point in the asset's desired resolution (scale equal to decimal display).
//
// Given the amount of asset units (U), and the number of units per BTC (Y), we
// compute the total amount of mSAT (X) as follows:
//   - X = (U / Y) * M
//   - where M is the number of mSAT in a BTC (100,000,000,000).
func UnitsToMilliSatoshi[N Int[N]](assetUnits,
	unitsPerBtc FixedPoint[N]) lnwire.MilliSatoshi {

	// Before we do the computation, we'll scale everything up to our
	// arithmetic scale.
	assetUnits = assetUnits.ScaleTo(arithScale)
	unitsPerBtc = unitsPerBtc.ScaleTo(arithScale)

	// We have the number of units, and the number of units per BTC, so we
	// can arrive at the number of of BTC via: BTC = units / (units/BTC).
	amtBTC := assetUnits.Div(unitsPerBtc)

	// Now that we have the amount of BTC, we can map to mSat by
	// multiplying by the number of mSAT in a BTC.
	oneBtcInMilliSat := FixedPointFromUint64[N](
		uint64(btcutil.SatoshiPerBitcoin*1_000), arithScale,
	)

	amtMsat := amtBTC.Mul(oneBtcInMilliSat)

	// We did the computation in terms of the scaled integers, so no we'll
	// go back to a normal mSAT value scaling down to zero (no decimals)
	// along the way.
	return lnwire.MilliSatoshi(amtMsat.ScaleTo(0).ToUint64())
}
