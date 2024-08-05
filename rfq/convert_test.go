package rfq

import (
	"fmt"
	"math"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/stretchr/testify/require"
)

var (
	btcPricesCents = []uint64{
		1_000_00,
		3_456_78,
		5_000_00,

		10_000_00,
		20_000_00,
		34_567_89,
		50_000_00,
		50_702_12,

		100_000_00,
		345_678_90,
		500_000_00,

		1_000_000_00,
		3_456_789_01,
		5_000_000_00,

		10_000_000_00,
		34_567_890_12,
		50_000_000_00,
	}

	maxDecimalDisplay = 8

	invoiceAmountsMsat = []uint64{
		1,
		2,
		3,
		5,

		10,
		34,
		50,

		100,
		345,
		500,

		1_000,
		3_456,
		5_000,

		10_000,
		34_567,
		50_000,

		100_000,
		345_678,
		500_000,

		1_000_000,
		3_456_789,
		5_000_000,

		10_000_000,
		20_000_000,
		34_567_890,
		50_000_000,

		100_000_000,
		345_678_901,
		500_000_000,

		1_000_000_000,
		3_456_789_012,
		5_000_000_000,

		10_000_000_000,
		34_567_890_123,
		50_000_000_000,

		100_000_000_000,
		345_678_901_234,
		500_000_000_000,
	}
)

// TestFindDecimalDisplayBoundaries tests the maximum number of units that can
// be represented with a given decimal display, the smallest payable invoice
// amount, and the maximum MPP rounding error. The values are printed out on
// standard output.
func TestFindDecimalDisplayBoundaries(t *testing.T) {
	for _, btcPriceCents := range btcPricesCents {
		fmt.Printf("-------------\nBTC price: %d USD\n-------------\n",
			btcPriceCents/100)
		for decDisp := 2; decDisp <= maxDecimalDisplay; decDisp++ {
			unitsPerUsd := uint64(math.Pow10(decDisp))

			priceCents := FixedPoint[uint64]{
				Scale: 2,
				Value: btcPriceCents,
			}
			priceScaled := priceCents.ScaleTo(decDisp)

			numShards := float64(16)
			maxUnits, smallestAmount, mSatPerUnit := calcLimits(
				btcPriceCents, decDisp,
			)

			maxRoundMSat := uint64(mSatPerUnit * numShards)

			oneUsd := FixedPoint[uint64]{
				Value: 1,
				Scale: 0,
			}.ScaleTo(decDisp)

			mSatPerUsd := UnitsToMilliSatoshi(
				oneUsd.Value, priceScaled,
			)
			unitsPerSat := MilliSatoshiToUnits(1000, priceScaled)

			fmt.Printf("decimalDisplay: %d\t\t\t%d units = 1 USD, "+
				"1 BTC = %d units\n"+
				"Max issuable units:\t\t\tcan represent %d "+
				"BTC\n"+
				"Min payable invoice amount:\t%d mSAT\n"+
				"Max MPP rounding error:\t\t%d mSAT (@%.0f "+
				"shards)\n"+
				"Satoshi per USD:\t\t\t%d\n"+
				"Satoshi per Asset Unit: \t%.5f\n"+
				"Asset Units per Satoshi: \t%d\n"+
				"Price In Asset: \t\t\t%d\n"+
				"Price Out Asset: \t\t\t%d\n\n",
				decDisp, unitsPerUsd, priceScaled.Value,
				maxUnits, smallestAmount, maxRoundMSat,
				numShards, mSatPerUsd/1000, mSatPerUnit/1000,
				unitsPerSat, priceScaled.Value,
				uint64(btcutil.SatoshiPerBitcoin*1000))
		}
	}
}

// calcLimits calculates the maximum number of units that can be represented
// with a given decimal display, the smallest payable invoice amount, and the
// maximum MPP rounding error for a given BTC price in cents, decimal display
// value and number of MPP shards.
func calcLimits(btcPriceCent uint64, decDisplay int) (uint64, uint64, float64) {
	msatScale := 11

	// In the unit test, the price is always given as cents per BTC.
	priceCents := FixedPoint[uint64]{
		Scale: 2,
		Value: btcPriceCent,
	}

	// priceScaled is the number of units per USD at the given price.
	priceScaled := priceCents.ScaleTo(decDisplay)

	// priceScaledF is the same as priceScaled, but in float64 format.
	priceScaledF := float64(btcPriceCent) * math.Pow10(decDisplay-2)

	// mSatPerUnitF is the number of mSAT per asset unit at the given
	// price.
	mSatPerUnitF := math.Pow10(msatScale) / priceScaledF

	// maxUnits is the maximum number of BTC that can be represented with
	// assets given the decimal display.
	maxUnits := uint64(math.MaxUint64) / priceScaled.Value

	smallestAmount := uint64(0)
	for _, invoiceAmount := range invoiceAmountsMsat {
		unitsForInvoice := (invoiceAmount * priceScaled.Value) /
			uint64(math.Pow10(msatScale))

		if unitsForInvoice > 0 && smallestAmount == 0 {
			smallestAmount = invoiceAmount
		}
	}

	return maxUnits, smallestAmount, mSatPerUnitF
}

// TestScaleTo tests the ScaleTo method of the FixedPoint type.
func TestScaleTo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		in          FixedPoint[uint64]
		scaleTo     int
		expectedOut FixedPoint[uint64]
	}{
		{
			name: "scale from 0 to 12",
			in: FixedPoint[uint64]{
				Value: 1,
				Scale: 0,
			},
			scaleTo: 12,
			expectedOut: FixedPoint[uint64]{
				Value: 1_000_000_000_000,
				Scale: 12,
			},
		},
		{
			name: "scale from 0 to 4",
			in: FixedPoint[uint64]{
				Value: 9,
				Scale: 0,
			},
			scaleTo: 4,
			expectedOut: FixedPoint[uint64]{
				Value: 90_000,
				Scale: 4,
			},
		},
		{
			name: "scale from 2 to 4",
			in: FixedPoint[uint64]{
				Value: 123_456,
				Scale: 2,
			},
			scaleTo: 4,
			expectedOut: FixedPoint[uint64]{
				Value: 12_345_600,
				Scale: 4,
			},
		},
		{
			name: "scale from 4 to 2, no precision loss",
			in: FixedPoint[uint64]{
				Value: 12_345_600,
				Scale: 4,
			},
			scaleTo: 2,
			expectedOut: FixedPoint[uint64]{
				Value: 123_456,
				Scale: 2,
			},
		},
		{
			name: "scale from 6 to 2, with precision loss",
			in: FixedPoint[uint64]{
				Value: 12_345_600,
				Scale: 6,
			},
			scaleTo: 2,
			expectedOut: FixedPoint[uint64]{
				Value: 1_234,
				Scale: 2,
			},
		},
		{
			name: "scale from 6 to 2, with full loss of value",
			in: FixedPoint[uint64]{
				Value: 12,
				Scale: 6,
			},
			scaleTo: 2,
			expectedOut: FixedPoint[uint64]{
				Value: 0,
				Scale: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.in.ScaleTo(tc.scaleTo)
			require.Equal(t, tc.expectedOut, out)
		})
	}
}

// TestMilliSatoshiToUnits tests the MilliSatoshiToUnits function.
func TestMilliSatoshiToUnits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		invoiceAmount lnwire.MilliSatoshi
		price         FixedPoint[uint64]
		expectedUnits uint64
	}{
		{
			// 5k USD per BTC @ decimal display 2.
			invoiceAmount: 200_000,
			price: FixedPoint[uint64]{
				Value: 5_000_00,
				Scale: 2,
			},
			expectedUnits: 1,
		},
		{
			// 5k USD per BTC @ decimal display 6.
			invoiceAmount: 200_000,
			price: FixedPoint[uint64]{
				Value: 5_000_00,
				Scale: 2,
			}.ScaleTo(6),
			expectedUnits: 10_000,
		},
		{
			// 50k USD per BTC @ decimal display 6.
			invoiceAmount: 1_973,
			price: FixedPoint[uint64]{
				Value: 50_702_00,
				Scale: 2,
			}.ScaleTo(6),
			expectedUnits: 1000,
		},
		{
			// 50M USD per BTC @ decimal display 6.
			invoiceAmount: 123_456_789,
			price: FixedPoint[uint64]{
				Value: 50_702_000_00,
				Scale: 2,
			}.ScaleTo(6),
			expectedUnits: 62595061158,
		},
		{
			// 50k USD per BTC @ decimal display 6.
			invoiceAmount: 5_070,
			price: FixedPoint[uint64]{
				Value: 50_702_12,
				Scale: 2,
			}.ScaleTo(6),
			expectedUnits: 2_570,
		},
		{
			// 7.341M JPY per BTC @ decimal display 6.
			invoiceAmount: 5_000,
			price: FixedPoint[uint64]{
				Value: 7_341_847,
				Scale: 0,
			}.ScaleTo(6),
			expectedUnits: 367_092,
		},
		{
			// 7.341M JPY per BTC @ decimal display 2.
			invoiceAmount: 5_000,
			price: FixedPoint[uint64]{
				Value: 7_341_847,
				Scale: 0,
			}.ScaleTo(4),
			expectedUnits: 3_670,
		},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("milliSat=%d,price=%s", tc.invoiceAmount,
			tc.price.String())

		t.Run(name, func(t *testing.T) {
			units := MilliSatoshiToUnits(tc.invoiceAmount, tc.price)
			require.Equal(t, tc.expectedUnits, units)

			mSat := UnitsToMilliSatoshi(units, tc.price)

			diff := tc.invoiceAmount - mSat
			require.LessOrEqual(t, diff, uint64(2), "mSAT diff")
		})
	}
}

// TestUsdToJpy tests the conversion of USD to JPY using a BTC price in USD and
// a BTC price in JPY, both expressed as a FixedPoint.
func TestUsdToJpy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		usdPrice    FixedPoint[uint64]
		jpyPrice    FixedPoint[uint64]
		usdAmount   uint64
		expectedJpy uint64
	}{
		{
			name: "1 USD to JPY @ 2.840M JPY/BTC, 20k USD/BTC",
			usdPrice: FixedPoint[uint64]{
				Value: 20_000_00,
				Scale: 2,
			}.ScaleTo(6),
			jpyPrice: FixedPoint[uint64]{
				Value: 2_840_000,
				Scale: 0,
			}.ScaleTo(4),
			usdAmount:   1,
			expectedJpy: 142,
		},
		{
			name: "100 USD to JPY @ 7.341M JPY/BTC, 50'702 USD/BTC",
			usdPrice: FixedPoint[uint64]{
				Value: 50_702_12,
				Scale: 2,
			}.ScaleTo(6),
			jpyPrice: FixedPoint[uint64]{
				Value: 7_341_847,
				Scale: 0,
			}.ScaleTo(4),
			usdAmount:   100,
			expectedJpy: 14_480,
		},
		{
			name: "500 USD to JPY @ 142M JPY/BTC, 1M USD/BTC",
			usdPrice: FixedPoint[uint64]{
				Value: 1_000_000_00,
				Scale: 2,
			}.ScaleTo(6),
			jpyPrice: FixedPoint[uint64]{
				Value: 142_000_000,
				Scale: 0,
			}.ScaleTo(4),
			usdAmount:   500,
			expectedJpy: 71_000,
		},
	}

	for _, tc := range testCases {
		// Easy way to scale up the USD amount to 6 decimal display.
		dollarUnits := FixedPoint[uint64]{
			Value: tc.usdAmount,
			Scale: 0,
		}.ScaleTo(6)

		// Convert the USD to mSAT.
		hundredUsdAsMilliSatoshi := UnitsToMilliSatoshi(
			dollarUnits.Value, tc.usdPrice,
		)

		// Convert the mSAT to JPY.
		usdAmountAsJpy := MilliSatoshiToUnits(
			hundredUsdAsMilliSatoshi, tc.jpyPrice,
		)

		// Go from decimal display of 4 to 0 (full JPY).
		fullJpy := usdAmountAsJpy / 10_000

		require.Equal(t, tc.expectedJpy, fullJpy)

		oneUsd := FixedPoint[uint64]{
			Value: 1,
			Scale: 0,
		}.ScaleTo(6)
		oneJpy := FixedPoint[uint64]{
			Value: 1,
			Scale: 0,
		}.ScaleTo(4)

		_, _, mSatPerUsdUnit := calcLimits(
			tc.usdPrice.ScaleTo(2).Value, 6,
		)
		mSatPerUsd := UnitsToMilliSatoshi(oneUsd.Value, tc.usdPrice)
		usdUnitsPerSat := MilliSatoshiToUnits(1000, tc.usdPrice)

		_, _, mSatPerJpyUnit := calcLimits(
			tc.jpyPrice.ScaleTo(0).Value, 4,
		)
		mSatPerJpy := UnitsToMilliSatoshi(oneJpy.Value, tc.jpyPrice)
		jpyUnitsPerSat := MilliSatoshiToUnits(1000, tc.jpyPrice)

		fmt.Printf("Satoshi per USD:\t\t\t\t%d\n"+
			"Satoshi per USD Asset Unit: \t%.5f\n"+
			"USD Asset Units per Satoshi: \t%d\n"+
			"Satoshi per JPY:\t\t\t\t%d\n"+
			"Satoshi per JPY Asset Unit: \t%.5f\n"+
			"JPY Asset Units per Satoshi: \t%d\n"+
			"Price In Asset: \t\t\t\t%d\n"+
			"Price Out Asset: \t\t\t\t%d\n"+
			"%3d USD in JPY: \t\t\t\t%d\n\n",
			mSatPerUsd/1000, mSatPerUsdUnit/1000, usdUnitsPerSat,
			mSatPerJpy/1000, mSatPerJpyUnit/1000, jpyUnitsPerSat,
			tc.usdPrice.Value, tc.jpyPrice.Value,
			tc.usdAmount, fullJpy)

	}
}
