package node_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "tkestack.io/tke/test/util/env"
)

func TestE2EInstaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Node Suite")
}
