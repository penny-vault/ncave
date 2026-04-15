package ncave_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/penny-vault/ncave/ncave"
	"github.com/penny-vault/pvbt/asset"
	"github.com/penny-vault/pvbt/data"
	"github.com/penny-vault/pvbt/engine"
	"github.com/penny-vault/pvbt/portfolio"
	"github.com/penny-vault/pvbt/universe"
)

var testTickers = []string{
	"CCCC", "CSTE", "SEER", "XBIT", "ORMP",
	"ACTG", "AMRX", "DIT", "FOSL", "LCUT",
	"AAPL", "MSFT", "NVDA", "JPM",
}

var _ = Describe("NetCurrentAssetValue", func() {
	var (
		ctx       context.Context
		snap      *data.SnapshotProvider
		nyc       *time.Location
		startDate time.Time
		endDate   time.Time
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		nyc, err = time.LoadLocation("America/New_York")
		Expect(err).NotTo(HaveOccurred())

		snap, err = data.NewSnapshotProvider("testdata/snapshot.db")
		Expect(err).NotTo(HaveOccurred())

		startDate = time.Date(2023, 6, 1, 0, 0, 0, 0, nyc)
		endDate = time.Date(2025, 7, 1, 0, 0, 0, 0, nyc)
	})

	AfterEach(func() {
		if snap != nil {
			snap.Close()
		}
	})

	runBacktest := func() portfolio.Portfolio {
		strategy := &ncave.NetCurrentAssetValue{
			Threshold: 1.5,
			Universe:  universe.NewStatic(testTickers...),
		}
		acct := portfolio.New(
			portfolio.WithCash(100000, startDate),
			portfolio.WithAllMetrics(),
		)

		eng := engine.New(strategy,
			engine.WithDataProvider(snap),
			engine.WithAssetProvider(snap),
			engine.WithAccount(acct),
		)

		result, err := eng.Backtest(ctx, startDate, endDate)
		Expect(err).NotTo(HaveOccurred())
		return result
	}

	It("produces expected returns and risk metrics", func() {
		result := runBacktest()

		summary, err := result.Summary()
		Expect(err).NotTo(HaveOccurred())
		Expect(summary.TWRR).To(BeNumerically("~", 0.4174, 0.001))
		Expect(summary.MaxDrawdown).To(BeNumerically("~", -0.2263, 0.001))
		Expect(result.Value()).To(BeNumerically("~", 141740, 50))
	})

	It("rebalances only at end of June", func() {
		result := runBacktest()
		txns := result.Transactions()

		rebalanceDates := map[string]bool{}
		for _, t := range txns {
			if t.Type == asset.BuyTransaction || t.Type == asset.SellTransaction {
				d := t.Date.In(nyc)
				Expect(d.Month()).To(Equal(time.June), "all trades should occur in June")
				rebalanceDates[d.Format("2006-01-02")] = true
			}
		}

		Expect(rebalanceDates).To(HaveKey("2023-06-30"))
		Expect(rebalanceDates).To(HaveKey("2025-06-30"))
	})

	It("selects the expected tickers on 2023-06-30 rebalance", func() {
		result := runBacktest()

		picks := map[string]bool{}
		for _, t := range result.Transactions() {
			if t.Type == asset.BuyTransaction &&
				t.Date.In(nyc).Format("2006-01-02") == "2023-06-30" {
				picks[t.Asset.Ticker] = true
			}
		}

		// 10 NCAV names pass the screen; the remaining 20 slots (out of the
		// default min-holdings of 30) are filled by QQQ.
		expected := []string{"ACTG", "AMRX", "CCCC", "CSTE", "DIT", "FOSL", "LCUT", "ORMP", "SEER", "XBIT", "QQQ"}
		Expect(picks).To(HaveLen(len(expected)))
		for _, ticker := range expected {
			Expect(picks).To(HaveKey(ticker), "expected %s in 2023-06-30 rebalance", ticker)
		}
	})

	It("drops stocks that no longer pass the threshold on 2025-06-30", func() {
		result := runBacktest()

		picks := map[string]bool{}
		for _, t := range result.Transactions() {
			if t.Type == asset.BuyTransaction &&
				t.Date.In(nyc).Format("2006-01-02") == "2025-06-30" {
				picks[t.Asset.Ticker] = true
			}
		}

		// ACTG and AMRX no longer pass the threshold in Q1 2025.
		expected := []string{"CCCC", "CSTE", "DIT", "FOSL", "LCUT", "ORMP", "SEER", "XBIT"}
		Expect(picks).To(HaveLen(len(expected)))
		for _, ticker := range expected {
			Expect(picks).To(HaveKey(ticker), "expected %s in 2025-06-30 rebalance", ticker)
		}
	})

	It("excludes financial-sector and high-cap non-value stocks", func() {
		result := runBacktest()

		excluded := map[string]bool{
			"JPM":  true, // financial sector
			"AAPL": true, // NCAV/MV too low
			"MSFT": true, // NCAV/MV too low
			"NVDA": true, // NCAV/MV too low
		}

		for _, t := range result.Transactions() {
			if t.Type == asset.BuyTransaction {
				Expect(excluded).NotTo(HaveKey(t.Asset.Ticker),
					"%s should not have been purchased", t.Asset.Ticker)
			}
		}
	})
})
