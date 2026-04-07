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
	"sort"
	"time"

	"github.com/penny-vault/pvbt/asset"
	"github.com/penny-vault/pvbt/data"
	"github.com/penny-vault/pvbt/engine"
	"github.com/penny-vault/pvbt/portfolio"
	"github.com/rs/zerolog"
)

//go:embed README.md
var description string

// NetCurrentAssetValue implements the NCAV/MV effect strategy.
// It buys stocks whose net current asset value exceeds a threshold
// multiple of their market capitalization, rebalancing annually in June.
type NetCurrentAssetValue struct {
	IndexName string  `pvbt:"index"     desc:"Stock index universe to select from"      default:"SPX" suggest:"SPX=SPX|NDX=NDX"`
	Threshold float64 `pvbt:"threshold" desc:"Minimum NCAV/MV ratio to include a stock" default:"1.5"`
}

func (s *NetCurrentAssetValue) Name() string {
	return "Net Current Asset Value"
}

func (s *NetCurrentAssetValue) Setup(_ *engine.Engine) {}

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
	// Only rebalance in June (portfolio formed annually per Quantpedia).
	if eng.CurrentDate().Month() != time.June {
		return nil
	}

	log := zerolog.Ctx(ctx)

	indexUniverse := eng.IndexUniverse(s.IndexName)

	fundDF, err := indexUniverse.At(ctx, data.WorkingCapital, data.MarketCap)
	if err != nil {
		return fmt.Errorf("fetch fundamentals: %w", err)
	}

	assets := fundDF.AssetList()
	if len(assets) == 0 {
		return nil
	}

	// Compute NCAV/MV = WorkingCapital / MarketCap and select stocks above threshold.
	type candidate struct {
		stock asset.Asset
		ratio float64
	}

	var selected []candidate

	for _, stock := range assets {
		wc := fundDF.Value(stock, data.WorkingCapital)
		mc := fundDF.Value(stock, data.MarketCap)

		if math.IsNaN(wc) || math.IsNaN(mc) || mc <= 0 || wc <= 0 {
			continue
		}

		ratio := wc / mc
		if ratio > s.Threshold {
			selected = append(selected, candidate{stock: stock, ratio: ratio})

			log.Debug().
				Str("ticker", stock.Ticker).
				Float64("working_capital", wc).
				Float64("market_cap", mc).
				Float64("ncav_mv", ratio).
				Msg("stock passes NCAV/MV threshold")
		}
	}

	if len(selected) == 0 {
		log.Debug().
			Int("universe_size", len(assets)).
			Float64("threshold", s.Threshold).
			Msg("no stocks pass NCAV/MV threshold")

		return nil
	}

	// Sort descending by NCAV/MV ratio.
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].ratio > selected[j].ratio
	})

	// Equal weight across all selected stocks.
	weight := 1.0 / float64(len(selected))
	members := make(map[asset.Asset]float64, len(selected))

	for _, c := range selected {
		members[c.stock] = weight
	}

	justification := fmt.Sprintf("%d stocks with NCAV/MV > %.1f from %s", len(selected), s.Threshold, s.IndexName)

	batch.Annotate("universe-size", fmt.Sprintf("%d", len(assets)))
	batch.Annotate("selected-count", fmt.Sprintf("%d", len(selected)))
	batch.Annotate("justification", justification)

	allocation := portfolio.Allocation{
		Date:          eng.CurrentDate(),
		Members:       members,
		Justification: justification,
	}

	return batch.RebalanceTo(ctx, allocation)
}
