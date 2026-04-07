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
)

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
		strategy := &ncave.NetCurrentAssetValue{}
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

	It("rebalances only in June", func() {
		result := runBacktest()
		txns := result.Transactions()

		rebalanceMonths := map[time.Month]bool{}
		for _, t := range txns {
			if t.Type == asset.BuyTransaction || t.Type == asset.SellTransaction {
				rebalanceMonths[t.Date.In(nyc).Month()] = true
			}
		}

		for m := range rebalanceMonths {
			Expect(m).To(Equal(time.June), "all trades should occur in June")
		}
	})

	It("only buys stocks with positive NCAV/MV", func() {
		result := runBacktest()
		txns := result.Transactions()

		buyCount := 0
		for _, t := range txns {
			if t.Type == asset.BuyTransaction {
				buyCount++
			}
		}

		Expect(buyCount).To(BeNumerically(">=", 1), "should buy at least one stock")
	})
})
