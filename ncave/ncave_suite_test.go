package ncave_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNetCurrentAssetValue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Net Current Asset Value Suite")
}
