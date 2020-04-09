// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package main

import (
	"flag"
	"fmt"
	"github.com/greenpau/go-idrac-redfish-api/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	appName        = "go-idrac-redfish-api-client"
	appDescription = "iDRAC Redfish API Client"
	appDocPath     = "https://github.com/greenpau/go-idrac-redfish-api/"
	appVersion     = "[untracked]"
	gitBranch      string
	gitCommit      string
	buildUser      string // whoami
	buildDate      string // date -u
)

func main() {
	// Initialize API client
	cli := client.NewClient()

	// Initialize CLI arguments
	var logLevel string
	var isShowVersion bool

	var host string
	var proto string
	var authUser string
	var authPass string
	var apiOperation string
	var apiResource string
	var configFile string
	var port int
	var secure bool

	flag.StringVar(&configFile, "config", "redfish.yaml", "configuration file")
	flag.StringVar(&host, "host", "", "target hostname or ip address")
	flag.IntVar(&port, "port", 443, "target port")
	flag.StringVar(&proto, "proto", "https", "transport protocol, either https or http")
	flag.BoolVar(&secure, "secure", false, "validate certificates (default: false)")
	flag.StringVar(&authUser, "username", "", "username")
	flag.StringVar(&authPass, "password", "", "password")
	flag.StringVar(&apiOperation, "operation", "", "operation")
	flag.StringVar(&apiResource, "resource", "", "resource")

	flag.StringVar(&logLevel, "log.level", "info", "logging severity level")
	flag.BoolVar(&isShowVersion, "version", false, "version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s - %s\n\n", appName, appDescription)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", appName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOperations:\n")
		for opName, opDescr := range cli.GetOperations() {
			fmt.Fprintf(os.Stderr, "  - %s: %s\n", opName, opDescr)
		}
		fmt.Fprintf(os.Stderr, "\nDocumentation: %s\n\n", appDocPath)
	}
	flag.Parse()
	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s %s", appName, appVersion)
		if gitBranch != "" {
			fmt.Fprintf(os.Stdout, ", branch: %s", gitBranch)
		}
		if gitCommit != "" {
			fmt.Fprintf(os.Stdout, ", commit: %s", gitCommit)
		}
		if buildUser != "" && buildDate != "" {
			fmt.Fprintf(os.Stdout, ", build on %s by %s", buildDate, buildUser)
		}
		fmt.Fprint(os.Stdout, "\n")
		os.Exit(0)
	}
	if level, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(level)
	} else {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	// Determine configuration file name and extension
	if configFile == "" {
		configFile = "redfish.yaml"
	}
	configFileExt := filepath.Ext(configFile)
	if configFileExt == "" {
		log.Fatalf("--config specifies a file without an extension, e.g. .yaml or .json")
	}
	configName := strings.TrimSuffix(configFile, configFileExt)

	// Define configuration parsing settings
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("IDRAC_API")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath("$HOME/.redfish")
	viper.BindEnv("username")
	viper.BindEnv("password")

	// Get environment variables
	if authUser == "" {
		if v := viper.Get("username"); v != nil {
			authUser = viper.Get("username").(string)
			log.Debugf("Setting username '%s' via IDRAC_API_USERNAME environment variable", authUser)
		}
	}

	if authPass == "" {
		if v := viper.Get("password"); v != nil {
			authPass = viper.Get("password").(string)
			log.Debugf("Setting password *** via IDRAC_API_PASSWORD environment variable")
		}
	}

	// Read configuration file
	if err := viper.ReadInConfig(); err == nil {
		if authUser == "" {
			if v := viper.Get("username"); v != nil {
				authUser = v.(string)
				log.Debugf("Setting username '%s' via configuration file", authUser)
			}
		}
		if authPass == "" {
			if v := viper.Get("password"); v != nil {
				authPass = v.(string)
				log.Debugf("Setting password *** via configuration file")
			}
		}
	} else {
		if !strings.Contains(err.Error(), "Not Found in") {
			log.Fatalf("configuration file error %s", err)
		}
	}

	// Configure the API client
	if err := cli.SetHost(host); err != nil {
		log.Fatalf("--host error: %s", err)
	}
	if err := cli.SetPort(port); err != nil {
		log.Fatalf("--port error: %s", err)
	}
	if err := cli.SetProtocol(proto); err != nil {
		log.Fatalf("--proto error: %s", err)
	}
	if err := cli.SetUsername(authUser); err != nil {
		log.Fatalf("--username error: %s", err)
	}
	if err := cli.SetPassword(authPass); err != nil {
		log.Fatalf("--password error: %s", err)
	}
	if secure {
		if err := cli.SetSecure(); err != nil {
			log.Fatalf("--secure error: %s", err)
		}
	}

	if apiOperation == "" && apiResource == "" {
		log.Fatalf("either --operation or --resource argument is required")
	}

	if apiOperation != "" && apiResource != "" {
		log.Fatalf("the --operation or --resource arguments are mutually exclusive")
	}

	log.Debugf("Host: %s", host)
	log.Debugf("Port: %d", port)
	log.Debugf("Protocol: %s", proto)
	log.Debugf("Certificate validation: %t", secure)
	log.Debugf("Username: %s", authUser)

	timerStartTime := time.Now()

	if apiOperation != "" {
		switch apiOperation {
		case "info":
			info, err := cli.GetInfo()
			if err != nil {
				log.Fatalf("%s", err)
			}
			fmt.Fprintf(os.Stdout, "Host: %s\n", host)
			fmt.Fprintf(os.Stdout, "Product: %s\n", info.Product)
			fmt.Fprintf(os.Stdout, "Service Tag: %s\n", info.ServiceTag)
			fmt.Fprintf(os.Stdout, "Manager MAC Address: %s\n", info.ManagerMACAddress)
			fmt.Fprintf(os.Stdout, "Redfish API Version: %s\n", info.RedfishVersion)
		default:
			log.Fatalf("the --operation %s is unsupported", apiOperation)
		}
	} else {
		res, err := cli.GetResource(apiResource)
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", res)
	}

	log.Debugf("took %s", time.Since(timerStartTime))
}
