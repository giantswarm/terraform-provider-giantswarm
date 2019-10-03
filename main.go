package main

import (
	"github.com/giantswarm/terraform-provider-giantswarm/provider"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
