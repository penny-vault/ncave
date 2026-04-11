// Copyright 2021-2026
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ncave

import (
	"context"
	_ "embed"
	"fmt"
	"math"
	"time"

	"github.com/penny-vault/pvbt/asset"
	"github.com/penny-vault/pvbt/data"
	"github.com/penny-vault/pvbt/engine"
	"github.com/penny-vault/pvbt/portfolio"
	"github.com/penny-vault/pvbt/universe"
	"github.com/rs/zerolog"
)

//go:embed README.md
var description string

// NetCurrentAssetValue implements the NCAV/MV effect strategy.
// It buys stocks from the US tradable universe whose net current asset value
// exceeds a threshold multiple of their market capitalization. Financial-sector
// stocks are excluded. The portfolio is rebalanced annually in July with equal
// weighting, using Q1 fundamental data to ensure all filings are available.
type NetCurrentAssetValue struct {
	Threshold float64           `pvbt:"threshold" desc:"Minimum NCAV/MV ratio to include a stock" default:"1.5"`
	Universe  universe.Universe `pvbt:"universe" desc:"Comma-separated tickers to constrain the universe (default: us-tradable index)"`
}

func (s *NetCurrentAssetValue) Name() string {
	return "Net Current Asset Value"
}

func (s *NetCurrentAssetValue) Setup(eng *engine.Engine) {
	if s.Universe == nil {
		s.Universe = eng.IndexUniverse("us-tradable")
		return
	}

	// User-supplied static universe carries only bare tickers; resolve each one
	// to a full asset record (with composite_figi and sector) so FetchAt and
	// sector filtering work correctly.
	bare := s.Universe.Assets(time.Time{})
	resolved := make([]asset.Asset, 0, len(bare))
	for _, a := range bare {
		resolved = append(resolved, eng.Asset(a.Ticker))
	}
	s.Universe = eng.Universe(resolved...)
}

func (s *NetCurrentAssetValue) Describe() engine.StrategyDescription {
	return engine.StrategyDescription{
		ShortCode:   "ncave",
		Description: description,
		Source:      "https://quantpedia.com/strategies/net-current-asset-value-effect/",
		Version:     "1.0.0",
		VersionDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
		Schedule:    "@monthend",
		Benchmark:   "VFINX",
	}
}

func (s *NetCurrentAssetValue) Compute(ctx context.Context, eng *engine.Engine, _ portfolio.Portfolio, batch *portfolio.Batch) error {
	// Rebalance at end of June (portfolio held from July per Quantpedia).
	// Tradecron doesn't cleanly support "@monthend in month 6" as a single spec
	// (combining @monthend with a Month cron field fires twice per June), so we
	// use @monthend and filter for June here.
	if eng.CurrentDate().Month() != time.June {
		return nil
	}

	log := zerolog.Ctx(ctx)

	// Get current universe members and exclude financial-sector stocks.
	currentDate := eng.CurrentDate()
	allMembers := s.Universe.Assets(currentDate)

	var nonFinancial []asset.Asset
	for _, a := range allMembers {
		if a.Sector != asset.SectorFinancialServices {
			nonFinancial = append(nonFinancial, a)
		}
	}

	if len(nonFinancial) == 0 {
		return nil
	}

	// Use Q1 (March 31) fundamental data so all filings are available by end of June.
	q1Date := time.Date(currentDate.Year(), time.March, 31, 16, 0, 0, 0, currentDate.Location())

	fundDF, err := eng.FetchAt(ctx, nonFinancial, q1Date, []data.Metric{data.WorkingCapital, data.MarketCap})
	if err != nil {
		return fmt.Errorf("fetch fundamentals: %w", err)
	}

	allAssets := fundDF.AssetList()
	if len(allAssets) == 0 {
		return nil
	}

	// Compute NCAV/MV = WorkingCapital / MarketCap and select stocks above threshold.
	type candidate struct {
		stock          asset.Asset
		ratio          float64
		workingCapital float64
		marketCap      float64
	}

	var selected []candidate

	for _, stock := range allAssets {
		wc := fundDF.Value(stock, data.WorkingCapital)
		mc := fundDF.Value(stock, data.MarketCap)

		if math.IsNaN(wc) || math.IsNaN(mc) || mc <= 0 || wc == 0 {
			continue
		}

		ratio := wc / mc
		if ratio <= s.Threshold {
			continue
		}

		selected = append(selected, candidate{
			stock:          stock,
			ratio:          ratio,
			workingCapital: wc,
			marketCap:      mc,
		})

		log.Debug().
			Str("ticker", stock.Ticker).
			Float64("working_capital", wc).
			Float64("market_cap", mc).
			Float64("ncav_mv", ratio).
			Msg("stock passes NCAV/MV threshold")
	}

	if len(selected) == 0 {
		log.Debug().
			Int("universe_size", len(allAssets)).
			Float64("threshold", s.Threshold).
			Msg("no stocks pass NCAV/MV threshold")

		return nil
	}

	// Equal weight across all selected stocks.
	weight := 1.0 / float64(len(selected))
	members := make(map[asset.Asset]float64, len(selected))

	for _, c := range selected {
		members[c.stock] = weight
		batch.Annotate(fmt.Sprintf("ncav-mv:%s", c.stock.Ticker), fmt.Sprintf("%.4f", c.ratio))
		batch.Annotate(fmt.Sprintf("wc:%s", c.stock.Ticker), fmt.Sprintf("%.0f", c.workingCapital))
		batch.Annotate(fmt.Sprintf("mc:%s", c.stock.Ticker), fmt.Sprintf("%.0f", c.marketCap))
	}

	justification := fmt.Sprintf("%d stocks with NCAV/MV > %.1f from US tradable universe (ex-financials)",
		len(selected), s.Threshold)

	batch.Annotate("universe-size", fmt.Sprintf("%d", len(nonFinancial)))
	batch.Annotate("selected-count", fmt.Sprintf("%d", len(selected)))
	batch.Annotate("justification", justification)

	allocation := portfolio.Allocation{
		Date:          currentDate,
		Members:       members,
		Justification: justification,
	}

	return batch.RebalanceTo(ctx, allocation)
}
