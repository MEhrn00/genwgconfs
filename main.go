package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"regexp"

	"github.com/mehrn00/genwgconfs/generator"
	"github.com/mehrn00/genwgconfs/utils"
)

// Global variable containing the command line arguments plus defaults
var CommandLineArgs utils.Arguments = utils.Arguments{
	// Set the default output to stdout
	Output: "stdout",

	// Set the default endpoint to an empty string
	Endpoint: "",

	// Set the default port to 51820
	Port: 51820,

	// Set the default subnet for the peer configs
	Subnet: net.IPNet{
		IP:   net.ParseIP("10.0.0.0"),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	},

	// Don't set AllowedIPs to 0.0.0.0/0. Generate a split-tunnel config by default
	All: false,

	// Generate 1 peer config by default
	Peers: 1,

	// PersistentKeepalive value. 0 means no PersistentKeepalive
	Keepalive: 0,
}

// Entrypoint for the program. Parses the command line arguments into a global variable
// named `CommandLineArgs`.
func init() {
	// Command line flag for the subnet of IP addresses each peer will receive
	flag.Func("subnet", "Subnet for the Wireguard configs. (default 10.0.0.0/24)", func(s string) error {
		ip, subnet, err := net.ParseCIDR(s)
		if err != nil {
			return errors.New("Invalid subnet.")
		}

		if !ip.Equal(subnet.IP) {
			mask, _ := subnet.Mask.Size()
			return errors.New(fmt.Sprintf("Invalid subnet base. Must be %s/%d for the defined mask", subnet.IP.String(), mask))
		}

		CommandLineArgs.Subnet = *subnet
		return nil
	})

	// Command line flag for specifying the endpoint IP address
	flag.Func("endpoint", "Endpoint IP address/hostname of the Wireguard server. (required)", func(s string) error {
		re := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

		// Check if the endpoint is possibly an IP address
		if re.Match([]byte(s)) {
			ip := net.ParseIP(s)
			if ip == nil {
				return errors.New("Endpoint is not a valid IP address or domain name.")
			}

			CommandLineArgs.Endpoint = s
		} else {
			CommandLineArgs.Endpoint = s
		}

		return nil
	})

	// Command line flag for specifying the Wireguard server's listen port
	flag.IntVar(&CommandLineArgs.Port, "port", 58120, "Server/Endpoint port.")

	// Command line flag for the output location
	flag.StringVar(&CommandLineArgs.Output, "output", "stdout", "Output file base name or 'stdout' to print the configs to stdout")

	// Command line flag for specifying the number of peer configs to generate
	flag.IntVar(&CommandLineArgs.Peers, "peers", 1, "Number of peer configs to generate.")

	// Command line flag for specifying if the peer configs should route all traffic
	// through the server
	flag.BoolVar(&CommandLineArgs.All, "all", false, "Set AllowedIPs to forward all traffic through the server.")

	/// Command line flag for PersistentKeepalive value
	flag.IntVar(&CommandLineArgs.Keepalive, "pk", 0, "Optional PersistentKeepalive value for peers. (default 0. No PersistentKeepalive value is set)")

	// Parse the command line flags (will check for flag parsing errors)
	flag.Parse()

	// Check if the user entered an endpoint IP address and fail out if an endpoint was
	// not specified
	if len(CommandLineArgs.Endpoint) == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "Endpoint not specified.")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Check if the supplied subnet allows for the number of peers specified
	ones, bits := CommandLineArgs.Subnet.Mask.Size()
	ipNum := int(math.Pow(2.0, float64(bits-ones))) - 1
	if ipNum < CommandLineArgs.Peers {
		fmt.Fprintln(flag.CommandLine.Output(), "Number of peers exceeds number of IPs in specified subnet. Use a larger subnet.")
		os.Exit(22)
	}
}

// Main function of the program. Command line argument logic occurs in the `init`
// function.
func main() {

	// Generate client configs based on the command line arguments
	clientConfs := generator.GenerateClientConfigs(&CommandLineArgs)
	serverConfig := generator.GenerateServerConfig(&CommandLineArgs, clientConfs)

	// Add the server public key to each of the clients
	for i := range clientConfs {
		clientConfs[i].AddServerKey(serverConfig.ServPublicKey)
	}

	// Check if printing to stdout or to files
	if CommandLineArgs.Output == "stdout" {
		// Print the server config
		fmt.Println("# Server -------------------------------------------------------")
		fmt.Println(serverConfig.Toini())

		// Print the client configs
		for i, v := range clientConfs {
			fmt.Printf("# Client %d -----------------------------------------------------\n", i+1)
			fmt.Println(v.Toini())
		}
	} else {
		baseName := CommandLineArgs.Output
		// Write out the server config
		os.WriteFile(fmt.Sprintf("%s_server.conf", baseName), []byte(serverConfig.Toini()), 0600)

		// Write out the client configs
		for i, v := range clientConfs {
			os.WriteFile(fmt.Sprintf("%s_client%d.conf", baseName, i+1), []byte(v.Toini()), 0600)
		}
	}
}
