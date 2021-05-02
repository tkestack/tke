package installer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tkestack.io/tke/pkg/util/log"
)

func TestTKE_newExpansionDriver(t *testing.T) {
	tke := &TKE{
		log: log.WithName("tke-installer"),
	}
	err := tke.newExpansionDriver()
	assert.True(t, err == nil)
}
