package cluster_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "tkestack.io/tke/test/util/env"
)

func TestCluster(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping cluster suite in short mode")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster Suite")
}
