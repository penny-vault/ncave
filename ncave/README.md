# Net Current Asset Value Effect

The Net Current Asset Value (NCAV) strategy is based on Benjamin Graham's approach of buying stocks trading below their net current asset value.

## Rules

The investment universe is the US tradable stock universe -- liquid US common stocks with a market-cap floor and contiguous price history. Financial-sector stocks and companies with missing fundamental data are excluded. The portfolio is rebalanced at the end of June using prior-fiscal-year fundamentals frozen as of a March 31 formation date, then held from July.

1. On the last trading day of June of year _Y_, screen the US tradable universe excluding financial-sector stocks.
2. For each stock, pull Working Capital from its most recent `ARQ` (as-reported, quarterly) filing with `date_key = Dec 31, Y-1`, restricted to filings with `event_date <= March 31, Y`. For calendar-fiscal-year companies this is the prior-year 10-K; for off-calendar companies it is the interim filing aligned to the Dec 31 date key. Late filers (10-Ks filed April-June of _Y_) are deliberately excluded.
3. Pull market capitalization for each stock from daily data at the March 31, _Y_ close.
4. Drop stocks with missing or invalid inputs (NaN working capital or market cap, non-positive market cap, or zero working capital).
5. Compute NCAV/MV = Working Capital / Market Capitalization for each remaining stock.
6. Select all stocks with NCAV/MV greater than the threshold (default: 1.5).
7. If the screen returns fewer than `min-holdings` stocks, each selected stock receives a `1/min-holdings` weight and the remaining portfolio weight is allocated to the `regime-change-asset`. Otherwise the selected stocks are equally weighted.
8. Hold the portfolio for one year until the next June rebalance.

### Why March 31 as the formation date

The March 31 formation date is this implementation's choice, not a rule from the cited literature. Two reasons drove it:

1. **10-K availability.** Calendar-fiscal-year 10-Ks are due within 60-90 days of Dec 31, so by March 31 the prior year's annual report is publicly available for the large majority of the universe.
2. **Excluding late filers.** Capping `event_date` at March 31 (rather than at the June rebalance date) excludes Q4 10-Ks filed between April and June. In the 2010-2026 backtest this produced materially higher CAGR than including late filers (28.3% vs 27.0%).

## Assets Typically Held

A variable number of deeply undervalued small-cap and micro-cap stocks that trade below their liquidation value. When the screen is sparse the portfolio also holds the regime-change asset (QQQ by default), which is a growth-oriented index used to capture the prevailing market regime when deep value is out of favor.

## Why the regime-change asset

The NCAV screen's breadth is itself a regime signal. The screen returns many names after broad market selloffs, when deep value is abundant and mean-reversion rallies are likely; it returns very few names in expensive, growth-led markets, when a concentrated basket of the few remaining deep-value stocks tends to underperform. Filling unfilled slots with a factor-opposite asset keeps the right-tail winners that drive the strategy's long-term returns while avoiding the concentration losses that historically occurred in sparse-screen years (2011, 2015, 2017, 2018, and 2024 in particular).

## Presets

- `Classic` -- disables the floor and holds only the stocks that pass the screen. Closest to Graham's original formulation.
- `Diversified` -- uses SPY as the regime-change asset instead of QQQ for broader, less growth-concentrated market exposure.
- `Defensive` -- uses a higher floor (50) and IVW (S&P 500 Growth) as the regime-change asset, trading some CAGR for a shallower drawdown profile.

## References

- Graham, B. and Dodd, D. (1934). _Security Analysis_. McGraw-Hill. Origin of the net-current-asset-value rule.
- Xiao, Y. and Arnold, G. (2008). "Testing Benjamin Graham's Net Current Asset Value Strategy in London." <https://papers.ssrn.com/sol3/papers.cfm?abstract_id=966188>
- Lauterbach, B. and Vu, J. (n.d.). "Ben Graham's Net Current Asset Value Rule Revisited: The Size-Adjusted Returns." <https://faculty.biu.ac.il/~lauteb/data_794/Ben_Graham.pdf>
- Oxman, J., Mohanty, S., and Carlisle, T. (2011). "Deep Value Investing and Unexplained Returns." <https://papers.ssrn.com/sol3/papers.cfm?abstract_id=1928694>
- Damodaran, A. (2012). "Value Investing: Investing for Grown Ups?" <https://papers.ssrn.com/sol3/papers.cfm?abstract_id=2042657>
