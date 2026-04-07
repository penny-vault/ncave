# Net Current Asset Value Effect

The Net Current Asset Value (NCAV) strategy is based on Benjamin Graham's approach of buying stocks trading below their net current asset value. It was formalized in the academic paper by [Quantpedia](https://quantpedia.com/strategies/net-current-asset-value-effect/).

## Rules

This strategy selects stocks whose net current asset value per share exceeds 1.5x their market capitalization per share. The portfolio is rebalanced annually in June.

1. On the last trading day of June, compute the NCAV/MV ratio for all stocks in the index universe. NCAV/MV = Working Capital / Market Capitalization.
2. Select all stocks with NCAV/MV greater than 1.5.
3. Invest equally in all selected stocks.
4. Hold the portfolio for one year until the next June rebalance.

## Assets Typically Held

A variable number of deeply undervalued small-cap and micro-cap stocks that trade below their liquidation value.
