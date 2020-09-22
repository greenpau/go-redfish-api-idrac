// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/greenpau/go-redfish-api-idrac/pkg/client"
	"github.com/greenpau/versioned"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	app        *versioned.PackageManager
	appVersion string
	gitBranch  string
	gitCommit  string
	buildUser  string
	buildDate  string
)

func init() {
	app = versioned.NewPackageManager("go-redfish-api-idrac-client")
	app.Description = "iDRAC Redfish API Client"
	app.Documentation = "https://github.com/greenpau/go-redfish-api-idrac/"
	app.SetVersion(appVersion, "1.0.1")
	app.SetGitBranch(gitBranch, "main")
	app.SetGitCommit(gitCommit, "69dfd98")
	app.SetBuildUser(buildUser, "")
	app.SetBuildDate(buildDate, "")
}

func main() {
	// Initialize API client
	cli := client.NewClient()
	supportedOperations := cli.GetOperations()

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
	var validateServerCert bool

	flag.StringVar(&configFile, "config", "redfish.yaml", "configuration file")
	flag.StringVar(&host, "host", "", "target hostname or ip address")
	flag.IntVar(&port, "port", 443, "target port")
	flag.StringVar(&proto, "proto", "https", "transport protocol, either https or http")
	flag.BoolVar(&validateServerCert, "validate-server-cert", false, "Verify the status of the server certificate")
	flag.StringVar(&authUser, "username", "", "username")
	flag.StringVar(&authPass, "password", "", "password")
	flag.StringVar(&apiOperation, "operation", "", "operation")
	flag.StringVar(&apiResource, "resource", "", "resource")

	flag.StringVar(&logLevel, "log.level", "info", "logging severity level")
	flag.BoolVar(&isShowVersion, "version", false, "version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s - %s\n\n", app.Name, app.Description)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", app.Name)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOperations:\n")
		for opName, op := range supportedOperations {
			fmt.Fprintf(os.Stderr, "  - %s: %s\n", opName, op.Description)
		}
		fmt.Fprintf(os.Stderr, "\nDocumentation: %s\n\n", app.Documentation)
	}
	flag.Parse()
	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s\n", app.Banner())
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
	viper.BindEnv("host")

	// Get environment variables

	if host == "" {
		if v := viper.Get("host"); v != nil {
			host = viper.Get("host").(string)
			log.Debugf("Setting host '%s' via IDRAC_API_HOST environment variable")
		}
	}

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
	if validateServerCert {
		if err := cli.SetValidateServerCertificate(); err != nil {
			log.Fatalf("--validate-server-cert error: %s", err)
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
	log.Debugf("Certificate validation: %t", validateServerCert)
	log.Debugf("Username: %s", authUser)

	timerStartTime := time.Now()

	if apiOperation != "" {
		if _, exists := supportedOperations[apiOperation]; !exists {
			log.Fatalf("the --operation %s is unsupported", apiOperation)
		}
		switch apiOperation {
		case "get-info":
			info, err := cli.GetInfo()
			if err != nil {
				log.Fatalf("%s", err)
			}
			fmt.Fprintf(os.Stdout, "Host: %s\n", host)
			fmt.Fprintf(os.Stdout, "Product: %s\n", info.Product)
			fmt.Fprintf(os.Stdout, "Service Tag: %s\n", info.ServiceTag)
			fmt.Fprintf(os.Stdout, "Manager MAC Address: %s\n", info.ManagerMACAddress)
			fmt.Fprintf(os.Stdout, "Redfish API Version: %s\n", info.RedfishVersion)
		case "get-systems":
			computerSystems, err := cli.GetComputerSystems()
			if err != nil {
				log.Fatalf("%s", err)
			}
			fmt.Fprintf(os.Stdout, "Number of Computer Systems: %d\n", len(computerSystems))
			fmt.Fprintf(os.Stdout, "---------------------------------\n")
			for _, cs := range computerSystems {
				if cs.Manufacturer != "" {
					fmt.Fprintf(os.Stdout, "System: %s | Manufacturer: %s\n", cs.ID, cs.Manufacturer)
				}
				if cs.Model != "" {
					fmt.Fprintf(os.Stdout, "System: %s | Model: %s\n", cs.ID, cs.Model)
				}
				if cs.SKU != "" {
					fmt.Fprintf(os.Stdout, "System: %s | SKU: %s\n", cs.ID, cs.SKU)
				}
				if cs.PartNumber != "" {
					fmt.Fprintf(os.Stdout, "System: %s | SKU: %s\n", cs.ID, cs.PartNumber)
				}
				if cs.BiosVersion != "" {
					fmt.Fprintf(os.Stdout, "System: %s | BIOS Version: %s\n", cs.ID, cs.BiosVersion)
				}
				spew.Dump(cs)
			}
		default:
			log.Fatalf("the --operation %s is supported by API, but not this utility", apiOperation)
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
