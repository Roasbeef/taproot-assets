# Asset Decimal Display

Within the Taproot Assets Protocol, an asset's unit is an integer (`uint64`).
That means, the protocol cannot represent fractions of an asset.
Therefore, any asset that represents a fiat currency would need to issue assets
equivalent to at least the smallest unit in use.
For example, an asset that represents the US-Dollar would need to be minted in a
way that one asset unit would represent one USD **cent**.
Or in other words, 100 units of such an asset would represent 1 US-Dollar.

Because wallet software user interfaces aren't expected to know what
"resolution" or precision any asset in the wild represents, a new field called
`decimal_display` was added to the JSON metadata of new assets minted with
`tapd v0.4.0-alpha` and later (this field is encoded in the metadata field as a
JSON field, therefore it is only compatible with assets that have a JSON
metadata field).

An issuer can specify `tapcli assets mint --decimal_display X` to specify the
number of decimal places the comma should be shifted to the left when displaying
a sum of asset units.

For the example above, a USD asset would choose `--decimal_display 2` to
indicate that 100 units ($10^2$) should be displayed as `$1.00` in the wallet
UI. Or another example: 1234 USD cent asset units would be displayed as `$12.34`
in the UI.

## Precision requirement for assets in the Lightning Network

Due to the non-divisible (integer) nature of Taproot Asset units, the smallest
asset amount that can be transferred with a Lightning Network HTLC is one asset
unit.

If one such asset unit represents significantly more than one satoshi, then in
some cases, due to integer division and rounding up, the user might end up
spending more assets than necessary when paying an invoice.

**Example 1: Paying a 1 satoshi invoice**:

While writing this article, one USD cent is roughly equivalent to 19 satoshi.

So if a user with USD cent assets in their wallet attempted to pay an invoice
denominated over 1 satoshi, they would need to send a full USD cent to satisfy
the invoice (again, only full asset units can be transported over an HTLC).

Even if one cent isn't much, the overpayment would still be roughly `19x`.

**Example 2: Paying an invoice with MPP**:

Multi-Part Payments (MPP) allow a single payment to be split up into multiple
parts/paths, resulting in multiple HTLCs per payment.

If a user wants to pay an invoice over `1,000,000` satoshi, that would be
equivalent to `52,631.5789` USD cents.
If the user were to pay this amount in a single HTLC, they would send `52,632`
USD cents (need to round up to satisfy integer asset amount and invoice minimum
amount requirements), overpaying by `0.4211` cents.

If the user's wallet decided to split up the payment into 16 parts for example,
then each part would correspond to `3289.4737` cents. To satisfy the integer
asset amount and invoice minimum amount requirement, each of the 16 HTLCs would
send out `3290` cents. That's a full `8.4211` cents of overpayment.

## What precision should I choose when minting an asset for the Lightning Network?

To address the issue of rounding up when splitting payments or representing
small satoshi amounts as asset units, an issuer of assets should use a high
enough value for `decimal_display` when minting.

**But what is a good value for `decimal_display`?**

We recommend to use a `decimal_display` value of `6` for currencies which
use a smaller subunit with two decimal places (such as cents for USD or EUR,
penny for GBP and so on).

For currencies without smaller units (for example JPY or VND), a
`decimal_display` value of `4` is recommended.

# RFQ

The RFQ system is responsible for acquiring real-time price quotes for
converting between asset units and satoshi (both directions) or between
different asset types.

There are two main user stories (as seen from the point of view of the wallet
end user):

1. Sending out assets: The user wants to pay a Lightning Network invoice that
   is denominated in satoshi. The user only has assets in their wallet, they
   (or their wallet software) want to find out how many asset units they need to
   send in order to satisfy the invoice amount in satoshi.
2. Receiving assets: The user wants to get paid in a specific asset. The user
   only knows about the asset, so they (or their wallet software) want to find
   out what the asset amount corresponds to in satoshi.

## Sell order

The sell order covers the first user story: The user wants to pay a
satoshi-denominated invoice with assets. The end result is that the user sells
some of their assets in their channel to the RFQ peer (edge node), requesting
them to forward satoshi to the network.

Formal definition (TODO: refine):
- Use case: sending assets as a payment, selling `buxx` for `msat`
- User query: `Q = how many out_asset units for in_asset amount?` (how many
  `buxx` do I need to seel to pay this payment denominated in `msat`?)
- `out_asset`: `buxx` (user sells `buxx` asset to RFQ peer)
- `in_asset`: `msat` (user "receives" `msat` from RFQ peer, which are then
  routed to the network)
- `max_amount`: `in_asset` (what is the maximum amount of `msat` the RFQ peer
  has to forward to the network? Equal to invoice amount plus user-defined max
  routing fee limit)
- `price_out_asset`: `out_asset_units_per_btc` (`buxx per BTC`)
- `price_in_asset`: `in_asset_units_per_btc` (`msat per BTC`)

## Buy order

The buy order covers the second user story: The user wants to get paid, they
create a satoshi-denominated invoice from their chosen asset amount they want to
receive. The end result is that the user buys assets into their asset channel
from their RFQ peer (edge node) and the RFQ peer is paid in satoshi by the
network.

Formal definition (TODO: refine):
- Use case: receiving assets through an invoice, selling `msat` for `buxx`
- User query: `Q = how many out_asset units for in_asset amount?` (how many
  `msat` should I denominate my invoice with to receive a given amount of
  `buxx`?)
- `out_asset`: `msat` (user sells sats to RFQ peer, which are routed to them by
  the network)
- `in_asset`: `buxx` (user buys `buxx` from RFQ peer)
- `max_amount`: `in_asset` (what is the maximum number of `buxx` the RFQ peer
  has to sell? Equal to the amount in the user query)
- `price_out_asset`: `out_asset_units_per_btc` (`buxx per BTC`)
- `price_in_asset`: `in_asset_units_per_btc` (`msat per BTC`)

## Examples

See `TestFindDecimalDisplayBoundaries` and `TestUsdToJpy`  in 
`rfq/convert_test.go` for how these examples are constructed.

**Case 1**: Buying/selling USD against BTC.

```text
In Asset:       USD with decimal display = 6 (1_000_000 asset units = 1 USD)
Out Asset:      satoshi / milli-satoshi

Example 1:
----------

What is price rate when 1 BTC = 20,000.00 USD?

decimalDisplay: 6			1000000 units = 1 USD, 1 BTC = 20000000000 units
Max issuable units:			can represent 922337203 BTC
Min payable invoice amount:	5 mSAT
Max MPP rounding error:		80 mSAT (@16 shards)
Satoshi per USD:			5000
Satoshi per Asset Unit: 	0.00500
Asset Units per Satoshi: 	200
Price In Asset: 			20000000000
Price Out Asset: 			100000000000


Example 2:
----------

What is price rate when 1 BTC = 1,000,000.00 USD?

decimalDisplay: 6			1000000 units = 1 USD, 1 BTC = 1000000000000 units
Max issuable units:			can represent 18446744 BTC
Min payable invoice amount:	1 mSAT
Max MPP rounding error:		1 mSAT (@16 shards)
Satoshi per USD:			100
Satoshi per Asset Unit: 	0.00010
Asset Units per Satoshi: 	10000
Price In Asset: 			1000000000000
Price Out Asset: 			100000000000


Example 3:
----------

What is price rate when 1 BTC = 10,000,000.00 USD?

decimalDisplay: 6			1000000 units = 1 USD, 1 BTC = 10000000000000 units
Max issuable units:			can represent 1844674 BTC
Min payable invoice amount:	1 mSAT
Max MPP rounding error:		0 mSAT (@16 shards)
Satoshi per USD:			10
Satoshi per Asset Unit: 	0.00001
Asset Units per Satoshi: 	100000
Price In Asset: 			10000000000000
Price Out Asset: 			100000000000
```

**Case 2**: Buying/selling USD against JPY.

```text
In Asset:       USD with decimal display = 6 (1_000_000 asset units = 1 USD)
Out Asset:      JPY with decimal display = 4 (10_000 asset units = 1 JPY)

Assumption:     1 USD = 142 JPY

Example 1:
----------

What is price rate when 1 BTC = 20,000.00 USD (1 BTC = 2,840,000 JPY)?

Satoshi per USD:				5000
Satoshi per USD Asset Unit: 	0.00500
USD Asset Units per Satoshi: 	200
Satoshi per JPY:				35
Satoshi per JPY Asset Unit: 	0.35211
JPY Asset Units per Satoshi: 	284
Price In Asset: 				20000000000
Price Out Asset: 				28400000000
  1 USD in JPY: 				142


Example 2:
----------

What is price rate when 1 BTC = 1,000,000.00 USD (1 BTC = 142,000,000 JPY)?

Satoshi per USD:				100
Satoshi per USD Asset Unit: 	0.00010
USD Asset Units per Satoshi: 	10000
Satoshi per JPY:				0
Satoshi per JPY Asset Unit: 	0.00704
JPY Asset Units per Satoshi: 	14199
Price In Asset: 				1000000000000
Price Out Asset: 				1420000000000
500 USD in JPY: 				71000
```
