package rates

import (
	"context"
	"fmt"
	"github.com/reddec/integration-lib/internal/utils"
	"github.com/shopspring/decimal"
	"net/url"
	"strings"
)

// Get major fiat currencies exchange rate to USD (alias to Fiat(USD))
func USD() (map[string]decimal.Decimal, error) {
	return Fiat("USD")
}

// Get exchange rate for fiat currency (see FiatContext for details)
func Fiat(base string) (map[string]decimal.Decimal, error) {
	return FiatContext(context.Background(), base)
}

// Get exchange rates for fiat currency using selected base (case insensitive) currency.
// Returned map will contain currency (in upper case) name and price in it for base currency.
//
// For example, for base currency USD, it will return JPY (Japan Yen)
// with value around 100: it means 1 USD is equal to 100 JPY.
//
// API from https://exchangeratesapi.io
func FiatContext(ctx context.Context, base string) (map[string]decimal.Decimal, error) {
	var reply struct {
		Rates map[string]decimal.Decimal `json:"rates"`
	}
	base = strings.TrimSpace(strings.ToUpper(base))
	var baseURL = "https://api.exchangeratesapi.io/latest?base=" + url.QueryEscape(base)
	err := utils.GetJSON(ctx, baseURL, &reply)
	if err != nil {
		return nil, fmt.Errorf("fiat rates: %w", err)
	}
	return reply.Rates, nil
}

// Get exchange rates for crypto currency quoted in USD (see CryptoContext for details)
func Crypto() (map[string]decimal.Decimal, error) {
	return CryptoContext(context.Background())
}

// Get exchange rates for crypto currency quoted in USD.
// Returned map will contain currency (in upper case) name and price in USD.
//
// For example, it will return BTC (Bitcoin) with value around 8000: it means 1 BTC is equal to $8000.
//
// API from https://coincap.io/
func CryptoContext(ctx context.Context) (map[string]decimal.Decimal, error) {
	const baseURL = `https://api.coincap.io/v2/rates`
	var reply struct {
		Data []struct {
			Symbol  string          `json:"symbol"`
			Type    string          `json:"type"`
			RateUSD decimal.Decimal `json:"rateUsd"`
		} `json:"data"`
	}
	err := utils.GetJSON(ctx, baseURL, &reply)
	if err != nil {
		return nil, fmt.Errorf("crypto rates: %w", err)
	}
	var ans = make(map[string]decimal.Decimal, len(reply.Data))
	for _, entry := range reply.Data {
		if entry.Type == "crypto" {
			ans[entry.Symbol] = entry.RateUSD
		}
	}
	return ans, nil
}
