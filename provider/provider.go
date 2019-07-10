package provider

import (
	"github.com/giantswarm/gsctl/client"
	"github.com/giantswarm/microerror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GIANTSWARM_INSTALLATION_ADDRESS", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GIANTSWARM_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"giantswarm_cluster": resourceGiantswarmCluster(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	address := d.Get("address").(string)
	token := d.Get("token").(string)

	clientV2, err := client.NewWithConfig(address, token)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return clientV2, nil
}
