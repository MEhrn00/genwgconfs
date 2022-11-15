package generator

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/mehrn00/genwgconfs/utils"
)

type ClientConfig struct {
	CliPublicKey        string
	ServPublicKey       string
	PrivateKey          string
	PresharedKey        string
	AllowedIPs          []net.IPNet
	Address             net.IPNet
	PersistentKeepAlive int
	Endpoint            string
}

func (c *ClientConfig) Toini() string {
	output := "[Interface]\n"
	output += fmt.Sprintf("Address = %s\n", c.Address.String())
	output += fmt.Sprintf("PrivateKey = %s\n\n", c.PrivateKey)
	output += "[Peer]\n"
	output += fmt.Sprintf("Endpoint = %s\n", c.Endpoint)
	output += fmt.Sprintf("PublicKey = %s\n", c.ServPublicKey)
	output += fmt.Sprintf("PresharedKey = %s\n", c.PresharedKey)

	if len(c.AllowedIPs) > 1 {
		output += "AllowedIPs = "
		for _, v := range c.AllowedIPs {
			output += v.String() + ","
		}

		output = strings.TrimSuffix(output, ",")
		output += "\n"
	} else {
		output += fmt.Sprintf("AllowedIPs = %s\n", c.AllowedIPs[0].String())
	}

	if c.PersistentKeepAlive > 0 {
		output += fmt.Sprintf("PersistentKeepalive = %d\n", c.PersistentKeepAlive)
	}

	return output
}

func (c *ClientConfig) AddServerKey(serverKey string) {
	c.ServPublicKey = serverKey
}

func ipv4addrToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func intToIpv4Addr(addr uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, addr)
	return ip
}

func GenerateClientConfigs(cliargs *utils.Arguments) []ClientConfig {
	clientConfigs := make([]ClientConfig, cliargs.Peers)
	assignIp := ipv4addrToInt(cliargs.Subnet.IP) + 1

	for i := range clientConfigs {
		clientConfigs[i].Endpoint = fmt.Sprintf("%s:%d", cliargs.Endpoint, cliargs.Port)

		privateKey, publicKey, err := utils.NewX25519pair()
		if err != nil {
			panic(err)
		}

		clientConfigs[i].PrivateKey = base64.StdEncoding.EncodeToString(privateKey)
		clientConfigs[i].CliPublicKey = base64.StdEncoding.EncodeToString(publicKey)

		psk, err := utils.Newx25519Key()
		if err != nil {
			panic(err)
		}

		clientConfigs[i].PresharedKey = base64.StdEncoding.EncodeToString(psk)
		clientConfigs[i].PersistentKeepAlive = cliargs.Keepalive

		address := net.IPNet{
			IP:   intToIpv4Addr(assignIp),
			Mask: cliargs.Subnet.Mask,
		}
		clientConfigs[i].Address = address
		assignIp++

		if cliargs.All {
			clientConfigs[i].AllowedIPs = []net.IPNet{
				{
					IP:   net.ParseIP("0.0.0.0"),
					Mask: net.IPv4Mask(0, 0, 0, 0),
				},
			}
		} else {
			clientConfigs[i].AllowedIPs = []net.IPNet{
				cliargs.Subnet,
			}
		}
	}

	return clientConfigs
}
