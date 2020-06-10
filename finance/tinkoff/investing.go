package tinkoff

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/reddec/integration-lib/internal/utils"
	"github.com/shopspring/decimal"
)

// Construct Tinkoff Invest API client using env variables.
//
// TINKOFF_INVEST_TOKEN - is mandatory variable, will panic if empty
func DefaultInvest() *Invest {
	cfg := &Invest{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	if cfg.Token == "" {
		panic("invest token required")
	}
	return cfg
}

// Configuration for the tinkof investing Open-API
//
// See https://tinkoffcreditsystems.github.io/invest-openapi/
type Invest struct {
	Token string `env:"TINKOFF_INVEST_TOKEN"` // Open-API token
}

// Get portfolio aggregated by currency with background context
func (cfg Invest) Portfolio() (map[string]decimal.Decimal, error) {
	return cfg.PortfolioContext(context.Background())
}

// Get portfolio aggregated by currency.
//
// API: https://api-invest.tinkoff.ru/
func (cfg Invest) PortfolioContext(ctx context.Context) (map[string]decimal.Decimal, error) {
	const baseURL = `https://api-invest.tinkoff.ru/openapi/portfolio`
	var reply struct {
		Payload struct {
			Positions []struct {
				InstrumentType string          `json:"instrumentType"`
				Ticker         string          `json:"ticker"`
				Balance        decimal.Decimal `json:"balance"`
				AveragePrice   struct {
					Value    decimal.Decimal `json:"value"`
					Currency string          `json:"currency"`
				} `json:"averagePositionPrice"`
				ExpectedYield struct {
					Value decimal.Decimal `json:"value"`
				} `json:"expectedYield"`
			}
		} `json:"payload"`
	}
	err := utils.GetJSONWithHeaders(ctx, baseURL, &reply, map[string]string{
		"Authorization": "Bearer " + cfg.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("tinkoff invest: %w", err)
	}

	var ans = make(map[string]decimal.Decimal, len(reply.Payload.Positions))
	for _, position := range reply.Payload.Positions {
		if position.InstrumentType == "Currency" {
			name := position.Ticker[:3]
			ans[name] = ans[name].Add(position.Balance)
		} else {
			name := position.AveragePrice.Currency
			amount := position.Balance.Mul(position.AveragePrice.Value)
			ans[name] = ans[name].Add(amount).Add(position.ExpectedYield.Value)
		}
	}
	return ans, nil
}
