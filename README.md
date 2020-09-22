# go-redfish-api-idrac

<a href="https://github.com/greenpau/go-redfish-api-idrac/actions/" target="_blank"><img src="https://github.com/greenpau/go-redfish-api-idrac/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/go-redfish-api-idrac" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>

iDRAC Redfish API client library written in Go.

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
  * [API Client](#api-client)
* [References](#references)

<!-- end-markdown-toc -->

## Getting Started

### API Client

Install the client by running:

```
go get -u github.com/greenpau/go-redfish-api-idrac/go-redfish-api-idrac-client
```

Prior to using the client, add your credentials via the following environment
variables:

```bash
export IDRAC_API_USERNAME=admin
export IDRAC_API_PASSWORD=secret
```

Additionally, there is an option for adding host via environment variables:

```bash
export IDRAC_API_HOST=10.10.10.10
```

Alternative, the credentials may be kept in `config.yaml` configuration file.
The binary searches for the file `$HOME/.redfish` directory.

Next, use the API in the following manner:

```bash
bin/go-redfish-api-idrac-client --host 10.10.10.10 --operation get-info --log.level debug
bin/go-redfish-api-idrac-client --host 10.10.10.10 --operation get-systems --log.level debug
```

The list of available operations (`--operation` argument) follows:

* `get-info`: Get basic information about a remote API endpoint
* `get-system`: Get system information

Additionally, the `--resource` argument accepts any valid Redfish API Endpoint:

```bash
go-redfish-api-idrac-client --host 10.10.10.10 --resource "/redfish/v1/Systems" --log.level debug
go-redfish-api-idrac-client --host 10.10.10.10 --resource "/redfish/v1/Systems/System.Embedded.1" --log.level debug
```

## References

* [Open Data Protocol (OData)](https://en.wikipedia.org/wiki/Open_Data_Protocol)
* [OData Version 4.01 Documentation](https://www.odata.org/documentation/)
