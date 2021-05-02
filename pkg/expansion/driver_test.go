package expansion

import (
	"encoding/json"
	"testing"
	"tkestack.io/tke/pkg/util/log"
)

func TestNewExpansionDriver(t *testing.T) {

	d, err := NewExpansionDriver(log.WithName("expansion-test"))
	if err != nil {
		t.Fatal(err)
	}
	d.printLayout()

}

func TestLayoutToTke(t *testing.T) {

	d, err := NewExpansionDriver(log.WithName("expansion-test"))
	if err != nil {
		t.Fatal(err)
	}
	tke := d.layout.genTKEConfig()
	b, err := json.MarshalIndent(tke, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tke.json %v", string(b))

}
