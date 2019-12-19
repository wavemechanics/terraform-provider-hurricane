This is a Terraform provider for the [Hurricane Electric DNS](https://dns.he.net) service.

Hurricane Electric supports updating existing A and AAAA records through an HTTP endpoint.
At present there is no API for creating or removing any records.
These actions must be done manually with the web interface.

But when there are pre-provisioned A records, this plugin can manage them.

# Build and install
```
$ git clone https://github.com/wavemechanics/terraform-provider-hurricane
$ cd terraform-provider-hurricane
$ make install
```
`make install` will compile the binary and copy it into ~/.terraform.d/plugins with a versioned name.

# Pre-provisioning records

Since there is no create/delete API, A records must be created with a specific IP address prior to running Terraform. The specific address to use is 127.250.202.222.
A record with this IP is "deleted" as far as this plugin is concerned.

# Terraform provider

Configure the provider like this:

```
provider "hurricane" {
    dns_endpoint = "https://dyn.dns.he.net/nic/update"
    dns_password = "some-password"
    dns_zone     = "your-hosted-domain.com"
}
```

The `dns_endpoint` shown is the default, so you won't usually need to specify it here.

`dns_password` is optional.
It is used if the resource block doesn't have a `password`.
The Hurricane setup allows you to set individual passwords per DNS name.

`dns_zone` is the DNS zone hosted by Hurricane.
Names in Terraform resource blocks are relative to this zone.
This is only used if the resource block doesn't have a `zone`.

# Resource records

Configure A records like this:

```
resource "hurricane_dns_rr" "myhost" {
    zone = "your-hosted-domain.com"
    name = "myhost"
    ip = "4.3.2.1"
    password = "specific-password"
}
```

`zone` overrides any default `dns_zone` in the provider.

`name` is the A record hostname. Required.

`ip` is the hostname's IP address. Required.

`password` overrides any default `dns_password` in the provider.

AAAA records are on the todo list.

