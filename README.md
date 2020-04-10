# go-idrac-redfish-api

![build](https://github.com/greenpau/go-idrac-redfish-api/workflows/core/badge.svg?branch=master)

iDRAC Redfish API client library written in Go.

## Table of Contents

* [Getting Started](#getting-started)

## Getting Started

## API Client

By running `make`, you will generate `bin/go-idrac-redfish-api-client` binary.

Prior to using the binary, add your credentials via the following environment
variables:

```bash
export IDRAC_API_USERNAME=admin
export IDRAC_API_PASSWORD=secret
```

Alternative, the credentials may be kept in `config.yaml` configuration file.
The binary searches for the file `$HOME/.redfish` directory.

Next, use the API in the following manner:

```
bin/go-idrac-redfish-api-client --host 10.10.10.10 --operation info --log.level debug
bin/go-idrac-redfish-api-client --host 10.10.10.10 --resource "/redfish/v1/" --log.level debug
```

The list of available operations (`--operation` argument) follows:

```
Operations:
  - info: Get basic information about a remote API endpoint
```

The `--resource` argument accepts any valid Redfish API Endpoint.
