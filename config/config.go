package config

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Config struct {
	Endpoint  string
	Password  string // if RR doesn't have its own
	Zone      string // if RR doesn't have its own
	Client    *http.Client
	UserAgent string
}

func Configure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	pluginVersion := strings.TrimPrefix(Ident, "$Id: ")
	pluginVersion = strings.TrimSuffix(pluginVersion, " $")

	log.Printf("[INFO] version %s", pluginVersion)

	return &Config{
		Endpoint:  d.Get("dns_endpoint").(string),
		Password:  d.Get("dns_password").(string),
		Zone:      d.Get("dns_zone").(string),
		Client:    &http.Client{},
		UserAgent: fmt.Sprintf("Terraform/%s terraform-hurricane-plugin/%s", terraformVersion, pluginVersion),
	}, nil
}
