package rfq

import (
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
	"golang.org/x/exp/constraints"
)

// FixedPoint is used to represent fixed point arithmetic for currency related
// calculations. A fixed point consists of a value, and a scale. The value is
// the integer representation of the number. The scale is used to represent the
// fractional/decimal component.
type FixedPoint[T constraints.Unsigned] struct {
	// Value is the value of the FixedPoint integer.
	Value T

	// Scale is used to represent the fractional component. This always
	// represents a power of 10. Eg: a scale value of 2 (two decimal places)
	// maps to a multiplication by 100.
	Scale int
}

// String returns the string version of the fixed point value.
func (f FixedPoint[T]) String() string {
	value := float64(f.Value) / math.Pow10(f.Scale)
	return fmt.Sprintf("%.*f", f.Scale, value)
}

// ScaleTo returns a new FixedPoint that is scaled up or down to the given
// scale.
func (f FixedPoint[T]) ScaleTo(newScale int) FixedPoint[T] {
	scaleDiff := newScale - f.Scale
	multiplier := math.Pow10(scaleDiff)
	newValue := float64(f.Value) * multiplier

	return FixedPoint[T]{
		Value: T(newValue),
		Scale: newScale,
	}
}

// MilliSatoshiToUnits converts the given milli-satoshi amount to units using
// the given price in units per bitcoin as a fixed point in the asset's desired
// resolution (scale equal to decimal display).
func MilliSatoshiToUnits(milliSat lnwire.MilliSatoshi,
	unitsPerBtc FixedPoint[uint64]) uint64 {

	oneBtcInMilliSatF := toBigFloat(btcutil.SatoshiPerBitcoin * 1_000)
	priceUnitsPerBtcF := toBigFloat(unitsPerBtc.Value)

	milliSatPerUnitF := new(big.Float)
	milliSatPerUnitF.Quo(oneBtcInMilliSatF, priceUnitsPerBtcF)

	invoiceAmountF := toBigFloat(uint64(milliSat))

	units := new(big.Float)
	units.Quo(invoiceAmountF, milliSatPerUnitF)

	result, _ := units.Uint64()
	return result
}

// UnitsToMilliSatoshi converts the given number of asset units to a
// milli-satoshi amount, using the given price in units per bitcoin as a fixed
// point in the asset's desired resolution (scale equal to decimal display).
func UnitsToMilliSatoshi(assetUnits uint64,
	unitsPerBtc FixedPoint[uint64]) lnwire.MilliSatoshi {

	oneBtcInMilliSatF := toBigFloat(btcutil.SatoshiPerBitcoin * 1_000)
	priceUnitsPerBtcF := toBigFloat(unitsPerBtc.Value)

	milliSatPerUnitF := new(big.Float)
	milliSatPerUnitF.Quo(oneBtcInMilliSatF, priceUnitsPerBtcF)

	milliSatF := new(big.Float)
	milliSatF.Mul(milliSatPerUnitF, toBigFloat(assetUnits))

	result, _ := milliSatF.Uint64()
	return lnwire.MilliSatoshi(result)
}

// toBigFloat returns the given uint64 as a big.Float.
func toBigFloat(d uint64) *big.Float {
	return new(big.Float).SetInt64(int64(d))
}
