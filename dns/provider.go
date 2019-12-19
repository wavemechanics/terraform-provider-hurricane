package dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/wavemechanics/terraform-provider-hurricane/config"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dns_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HURRICANE_DNS_ENDPOINT", "https://dyn.dns.he.net/nic/update"),
				Description: "Hurricane Electric dynamic DNS URL for POST requests.",
			},

			"dns_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HURRICANE_DNS_PASSWORD", nil),
				Description: "Default password for dynamic DNS requests.",
			},

			"dns_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HURRICANE_DNS_ZONE", nil),
				Description: "Default Hurricane DNS zone.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"hurricane_dns_rr": resourceRR(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return config.Configure(d, terraformVersion)
	}
	return p
}
