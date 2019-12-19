package dns

import (
	"net"
	"net/http"
	"net/url"

	"github.com/bogdanovich/dns_resolver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/wavemechanics/etype"
	"github.com/wavemechanics/terraform-provider-hurricane/config"
)

const (
	// Pre-provisioned and "deleted" A records have the placeholder IP.
	// It is formed from 127 0xfa 0xca 0xde.
	PlaceholderIP = "127.250.202.222"
)

const (
	ErrAmbiguous     = etype.Sentinel("Ambiguous: multiple A records with the same name")
	ErrNoPlaceholder = etype.Sentinel("Existing A record must be placeholder")
	ErrUpdateFailed  = etype.Sentinel("Hurricane update call failed")
	ErrNXDOMAIN      = etype.Sentinel("Faked NXDOMAIN (really placeholder)")
)

func resourceRR() *schema.Resource {
	return &schema.Resource{
		Create: rrCreate,
		Read:   rrRead,
		Update: rrUpdate,
		Delete: rrDelete,

		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func rrCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*config.Config)

	zone := d.Get("zone").(string)
	if zone == "" {
		zone = provider.Zone
	}

	name := d.Get("name").(string)

	// RR must exist and have PlaceholderIP
	currentip, err := lookupIP(zone, name)
	if err != nil {
		return err
	}
	if currentip != PlaceholderIP {
		return ErrNoPlaceholder
	}

	ip := d.Get("ip").(string)
	err = updateIP(d, m, ip)
	if err != nil {
		return err
	}

	d.SetId(name + "." + zone)
	return rrRead(d, m)
}

func rrRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(*config.Config)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)

	if zone == "" {
		zone = provider.Zone
	}

	ip, err := lookupIP(zone, name)
	if err != nil {
		if err.Error() == "NXDOMAIN" {
			return ErrNoPlaceholder
		}
		return err
	}
	if ip != PlaceholderIP {
		d.Set("ip", ip)
	}
	return nil
}

func rrUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(*config.Config)

	zone := d.Get("zone").(string)
	if zone == "" {
		zone = provider.Zone
	}

	name := d.Get("name").(string)

	// RR must exist and must not have PlaceholderIP
	currentip, err := lookupIP(zone, name)
	if err != nil {
		return err
	}
	if currentip == PlaceholderIP {
		return ErrNXDOMAIN
	}

	ip := d.Get("ip").(string)
	err = updateIP(d, m, ip)
	if err != nil {
		return err
	}
	return rrRead(d, m)
}

func rrDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(*config.Config)

	zone := d.Get("zone").(string)
	if zone == "" {
		zone = provider.Zone
	}

	name := d.Get("name").(string)

	// RR must exist
	currentip, err := lookupIP(zone, name)
	if err != nil {
		return err
	}
	if currentip == PlaceholderIP {
		return nil // already "deleted"
	}

	err = updateIP(d, m, PlaceholderIP)
	if err != nil {
		return err
	}
	return nil
}

// lookupA looks up the current A record for the given hostname from
// the zone's authoritative nameservers.
// Must ask the authoritative nameservers directly to avoid caches.
//
// Returns error if there is more than one A record for the hostname.
// Since we are not really creating and deleting records, we fake it
// with the PlaceholderIP address.
// So it is an error if the record really doesn't exist (NXDOMAIN).
func lookupIP(zone, name string) (string, error) {
	nsrecs, err := lookupNS(zone)
	if err != nil {
		return "", err
	}

	resolver := dns_resolver.New(nsrecs)
	resolver.RetryTimes = 5

	ips, err := resolver.LookupHost(name + "." + zone)
	if err != nil {
		if err.Error() == "NXDOMAIN" {
			return "", ErrNoPlaceholder
		}
		return "", err
	}
	if len(ips) > 1 {
		return "", ErrAmbiguous
	}
	return ips[0].String(), nil
}

// lookupNS looks up a zone's authoritative nameservers.
func lookupNS(zone string) ([]string, error) {
	NS, err := net.LookupNS(zone)
	if err != nil {
		return nil, err
	}
	var hosts []string
	for _, ns := range NS {
		hosts = append(hosts, ns.Host)
	}
	return hosts, nil
}

// updateIP updates and existing A RR using the Hurricane HTTP endpoint.
func updateIP(d *schema.ResourceData, m interface{}, ip string) error {
	provider := m.(*config.Config)

	password := d.Get("password").(string)
	if password == "" {
		password = provider.Password
	}

	zone := d.Get("zone").(string)
	if zone == "" {
		zone = provider.Zone
	}

	name := d.Get("name").(string)

	// make POST request to update the existing A record
	//Authentication and Updating using a POST
	// % curl "https://dyn.dns.he.net/nic/update" -d "hostname=dyn.example.com" -d "password=password" -d "myip=192.168.0.1"
	// % curl "https://dyn.dns.he.net/nic/update" -d "hostname=dyn.example.com" -d "password=password" -d "myip=2001:db8:beef:cafe::1"

	values := url.Values{
		"hostname": {name + "." + zone},
		"password": {password},
		"myip":     {ip},
	}

	resp, err := http.PostForm(provider.Endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrUpdateFailed
	}

	// not sure if reading the body helps if we know we got a 200?
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//    return err
	//}
	//if string(body) != "good "+ip {
	//	return ErrUpdateFailed
	//}

	return nil
}
