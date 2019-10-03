package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"giantswarm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GIANTSWARM_INSTALLATION_ADDRESS"); v == "" {
		t.Fatal("GIANTSWARM_INSTALLATION_ADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("GIANTSWARM_TOKEN"); v == "" {
		t.Fatal("GIANTSWARM_TOKEN must be set for acceptance tests")
	}
}
