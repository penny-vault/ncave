# Net Current Asset Value Effect

The Net Current Asset Value (NCAV) strategy is based on Benjamin Graham's approach of buying stocks trading below their net current asset value. It was formalized in the academic paper by [Quantpedia](https://quantpedia.com/strategies/net-current-asset-value-effect/).

## Rules

The investment universe is the US tradable stock universe -- liquid US common stocks with a market-cap floor and contiguous price history. Financial-sector stocks and companies with missing fundamental data are excluded. The portfolio is rebalanced at the end of June using Q1 fundamental data to ensure all filings are publicly available, then held from July.

1. On the last trading day of June, fetch Q1 (March 31) working capital and market capitalization for the current US tradable universe, excluding financial-sector stocks.
2. Drop stocks with missing or invalid fundamentals (NaN, non-positive market cap, or zero working capital).
3. Compute NCAV/MV = Working Capital / Market Capitalization for each remaining stock.
4. Select all stocks with NCAV/MV greater than the threshold (default: 1.5).
5. Invest equally in all selected stocks.
6. Hold the portfolio for one year until the next June rebalance.

## Assets Typically Held

A variable number of deeply undervalued small-cap and micro-cap stocks that trade below their liquidation value.
