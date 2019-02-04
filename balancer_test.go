package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

const unexpectedError string = "got an error but didn't want one"
const missingError string = "wanted an error but didn't get one"
const wrongError string = "got an error but expected a different one"

func TestErrInvalidAssetAmount_Error(t *testing.T) {
	asset := Asset("ETH")
	amount := decimal.NewFromFloat(-5)

	err := ErrInvalidAssetAmount{Asset: asset, Amount: amount}

	want := "ETH needs positive amount, not -5"
	got := err.Error()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestSetPricelist(t *testing.T) {
	t.Run("a new pricelist can be set and retrieved", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}
	})
	t.Run("an empty pricelist cannot be set", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{})

		if err == nil {
			t.Error(missingError)
		}

		if err != ErrEmptyPricelist {
			t.Error(wrongError)
		}
	})
	t.Run("pricelist asset keys must be uppercase", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"btc": decimal.NewFromFloat(5000),
		})

		if err == nil {
			t.Error(missingError)
		}

		if err != ErrInvalidAsset {
			t.Error(wrongError)
		}
	})
	t.Run("pricelist entries must have a value above 0", func(t *testing.T) {
		invalidAsset := Asset("BTC")
		invalidAmount := decimal.NewFromFloat(-5)

		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH":        decimal.NewFromFloat(200),
			invalidAsset: invalidAmount,
		})

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}

func TestNewHoldings(t *testing.T) {
	got, err := NewHoldings(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(5),
	})

	if err != nil {
		t.Error(unexpectedError)
	}

	want := Holdings{"ETH": decimal.NewFromFloat(5)}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestNewHoldings_ErrorsOnNonPositiveHoldingAmount(t *testing.T) {
	asset := Asset("ETH")
	amount := decimal.NewFromFloat(-5)

	_, err := NewHoldings(map[Asset]decimal.Decimal{
		asset: amount,
	})

	want := ErrInvalidAssetAmount{Asset: asset, Amount: amount}

	if err != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestNewHoldings_ErrorsOnInvalidInput(t *testing.T) {
	testCases := []struct {
		name     string
		holdings map[Asset]decimal.Decimal
		err      error
	}{
		{
			name:     "holdings must not be empty",
			holdings: map[Asset]decimal.Decimal{},
			err:      ErrEmptyHoldings,
		},
		{
			name: "holding assets should be uppercase and unique",
			holdings: map[Asset]decimal.Decimal{
				"eth": decimal.NewFromFloat(5),
			},
			err: ErrInvalidAsset,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHoldings(tt.holdings)

			if err == nil {
				t.Error(missingError)
			}

			if err != tt.err {
				t.Errorf("got %v want %v", err, tt.err)
			}
		})
	}
}

func TestNewAccount(t *testing.T) {
	t.Run("a new account cannot be created with negative holding values", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		holdings := Holdings{
			"ETH": decimal.NewFromFloat(-5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(holdings)

		if err == nil {
			t.Error(missingError)
		}
	})
	t.Run("a new account cannot be created with empty holdings", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		holdings := Holdings{}

		_, err := NewAccount(holdings)

		if err == nil {
			t.Error(missingError)
		}
	})
	t.Run("a new account cannot be created with invalid asset names in its holdings", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		holdings := Holdings{
			"eth": decimal.NewFromFloat(5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(holdings)

		if err == nil {
			t.Error(missingError)
		}
	})
	t.Run("a new account can be created", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		holdings := Holdings{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(holdings)

		if err != nil {
			t.Error(unexpectedError)
		}
	})
}

func TestNewIndex(t *testing.T) {
	t.Run("an index cannot be empty", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{})

		if err == nil {
			t.Error(missingError)
		}

		want := ErrEmptyIndex

		if err != want {
			t.Error(wrongError)
		}
	})
	t.Run("index asset keys must be uppercase", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"btc": decimal.NewFromFloat(5000),
		})

		if err == nil {
			t.Error(missingError)
		}

		want := ErrInvalidAsset

		if err != want {
			t.Error(wrongError)
		}
	})
	t.Run("index entries must have a value above 0", func(t *testing.T) {
		invalidAsset := Asset("BTC")
		invalidAmount := decimal.NewFromFloat(-5)

		_, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH":        decimal.NewFromFloat(200),
			invalidAsset: invalidAmount,
		})

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
	t.Run("index values must sum to 1", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.2),
			"BTC": decimal.NewFromFloat(0.2),
		})

		if err == nil {
			t.Error(missingError)
		}

		want := ErrIndexSumIncorrect

		if err != want {
			t.Error(wrongError)
		}
	})
	t.Run("a new index can be created", func(t *testing.T) {
		got, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.5),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		want := Index{
			"ETH": decimal.NewFromFloat(0.5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestAccount_Balance(t *testing.T) {
	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	holdings := Holdings{
		"ETH": decimal.NewFromFloat(20),
		"BTC": decimal.NewFromFloat(0.5),
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.3),
		"BTC": decimal.NewFromFloat(0.7),
	}

	Account, _ := NewAccount(holdings)

	got, err := Account.Balance(targetIndex)

	if err != nil {
		t.Error(unexpectedError)
	}

	want := map[Asset]Trade{
		"ETH": {Action: "sell", Amount: decimal.NewFromFloat(10.25)},
		"BTC": {Action: "buy", Amount: decimal.NewFromFloat(0.41)},
	}

	assertSameTrades(t, got, want)
}

func TestAccount_Balance_IntoNewAssets(t *testing.T) {
	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH":  decimal.NewFromFloat(200),
		"BTC":  decimal.NewFromFloat(2000),
		"IOTA": decimal.NewFromFloat(0.3),
		"BAT":  decimal.NewFromFloat(0.12),
		"XLM":  decimal.NewFromFloat(0.2),
	})

	holdings := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(42),
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH":  decimal.NewFromFloat(0.2),
		"BTC":  decimal.NewFromFloat(0.2),
		"IOTA": decimal.NewFromFloat(0.2),
		"BAT":  decimal.NewFromFloat(0.2),
		"XLM":  decimal.NewFromFloat(0.2),
	}

	Account, _ := NewAccount(holdings)

	got, err := Account.Balance(targetIndex)

	if err != nil {
		t.Error(unexpectedError)
	}

	want := map[Asset]Trade{
		"ETH":  {Action: "sell", Amount: decimal.NewFromFloat(33.6)},
		"BTC":  {Action: "buy", Amount: decimal.NewFromFloat(0.84)},
		"IOTA": {Action: "buy", Amount: decimal.NewFromFloat(5600)},
		"BAT":  {Action: "buy", Amount: decimal.NewFromFloat(14000)},
		"XLM":  {Action: "buy", Amount: decimal.NewFromFloat(8400)},
	}

	assertSameTrades(t, got, want)
}

func TestAccount_Balance_ErrorsWhenTargetIndexIsInvalid(t *testing.T) {
	testCases := []struct {
		name        string
		targetIndex map[Asset]decimal.Decimal
	}{
		{
			name: "target index does not sum to 1",
			targetIndex: map[Asset]decimal.Decimal{
				"ETH": decimal.NewFromFloat(0.2),
				"BTC": decimal.NewFromFloat(0.2),
			},
		},
		{
			name:        "target index is empty",
			targetIndex: map[Asset]decimal.Decimal{},
		},
		{
			name: "target index has an asset missing from the pricelist",
			targetIndex: map[Asset]decimal.Decimal{
				"ETH": decimal.NewFromFloat(0.8),
				"BAT": decimal.NewFromFloat(0.2),
			},
		},
	}

	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	holdings := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(42),
	}

	Account, _ := NewAccount(holdings)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Account.Balance(tt.targetIndex)

			if err == nil {
				t.Error(missingError)
			}
		})
	}
}

func assertSameTrades(t *testing.T, got map[Asset]Trade, want map[Asset]Trade) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got %d trades want %d", len(got), len(want))
	}

	for asset, wantTrade := range want {
		gotTrade, exists := got[asset]
		if !exists {
			t.Fatalf("asset %s missing from trade list", asset)
		}
		if gotTrade.Action != wantTrade.Action {
			t.Fatalf(
				"got a trade action of %s, want %s for asset %s",
				gotTrade.Action,
				wantTrade.Action,
				asset,
			)
		}
		if !gotTrade.Amount.Equal(wantTrade.Amount) {
			t.Fatalf(
				"got %v want %v for trade of asset %s",
				gotTrade.Amount,
				wantTrade.Amount,
				asset,
			)
		}
	}
}

func ExampleAccount_Balance() {
	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	holdings := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(20),
		"BTC": decimal.NewFromFloat(0.5),
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.5),
		"BTC": decimal.NewFromFloat(0.5),
	}

	Account, _ := NewAccount(holdings)

	requiredTrades, _ := Account.Balance(targetIndex)

	for asset, trade := range requiredTrades {
		fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
	}

	// Unordered output:
	// sell 3.75 ETH
	// buy 0.15 BTC
}

func BenchmarkBalance(b *testing.B) {
	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	for i := 0; i < b.N; i++ {
		holdings := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		}
		targetIndex := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		}

		Account, _ := NewAccount(holdings)

		_, _ = Account.Balance(targetIndex)
	}
}
