package e2e_installer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "tkestack.io/tke/test/util/env"
)

func TestE2EInstaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2EInstaller Suite")
}
