# Net Current Asset Value Effect

The Net Current Asset Value (NCAV) strategy is based on Benjamin Graham's approach of buying stocks trading below their net current asset value. It was formalized in the academic paper by [Quantpedia](https://quantpedia.com/strategies/net-current-asset-value-effect/).

## Rules

The investment universe is the US tradable stock universe -- liquid US common stocks with a market-cap floor and contiguous price history. Financial-sector stocks and companies with missing fundamental data are excluded. The portfolio is rebalanced at the end of June using Q1 fundamental data to ensure all filings are publicly available, then held from July.

1. On the last trading day of June, fetch Q1 (March 31) working capital and market capitalization for the current US tradable universe, excluding financial-sector stocks.
2. Drop stocks with missing or invalid fundamentals (NaN, non-positive market cap, or zero working capital).
3. Compute NCAV/MV = Working Capital / Market Capitalization for each remaining stock.
4. Select all stocks with NCAV/MV greater than the threshold (default: 1.5).
5. If the screen returns fewer than `min-holdings` stocks, each selected stock receives a `1/min-holdings` weight and the remaining portfolio weight is allocated to the `regime-change-asset`. Otherwise the selected stocks are equally weighted.
6. Hold the portfolio for one year until the next June rebalance.

## Assets Typically Held

A variable number of deeply undervalued small-cap and micro-cap stocks that trade below their liquidation value. When the screen is sparse the portfolio also holds the regime-change asset (QQQ by default), which is a growth-oriented index used to capture the prevailing market regime when deep value is out of favor.

## Why the regime-change asset

The NCAV screen's breadth is itself a regime signal. The screen returns many names after broad market selloffs, when deep value is abundant and mean-reversion rallies are likely; it returns very few names in expensive, growth-led markets, when a concentrated basket of the few remaining deep-value stocks tends to underperform. Filling unfilled slots with a factor-opposite asset keeps the right-tail winners that drive the strategy's long-term returns while avoiding the concentration losses that historically occurred in sparse-screen years (2011, 2015, 2017, 2018, and 2024 in particular).

## Presets

- `Classic` -- disables the floor and holds only the stocks that pass the screen. Closest to Graham's original formulation.
- `Diversified` -- uses SPY as the regime-change asset instead of QQQ for broader, less growth-concentrated market exposure.
- `Defensive` -- uses a higher floor (50) and IVW (S&P 500 Growth) as the regime-change asset, trading some CAGR for a shallower drawdown profile.
