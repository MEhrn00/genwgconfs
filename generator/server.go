package generator

import (
	"encoding/base64"
	"fmt"
	"net"

	"github.com/mehrn00/genwgconfs/utils"
)

type ServerConfig struct {
	address       string
	privateKey    string
	ServPublicKey string
	listenPort    int
	peers         []ServerPeer
}

type ServerPeer struct {
	allowedIPs   net.IPNet
	publicKey    string
	presharedKey string
}

func (p *ServerPeer) Toini() string {
	output := "[Peer]\n"
	output += fmt.Sprintf("AllowedIPs = %s/32\n", p.allowedIPs.IP.String())
	output += fmt.Sprintf("PublicKey = %s\n", p.publicKey)
	output += fmt.Sprintf("PresharedKey = %s\n", p.presharedKey)
	return output
}

func (s *ServerConfig) Toini() string {
	output := "[Interface]\n"
	output += fmt.Sprintf("Address = %s\n", s.address)
	output += fmt.Sprintf("PrivateKey = %s\n", s.privateKey)
	output += fmt.Sprintf("ListenPort = %d\n", s.listenPort)

	for _, v := range s.peers {
		output += "\n"
		output += v.Toini()
	}

	return output
}

func GenerateServerConfig(cliargs *utils.Arguments, clientConfigs []ClientConfig) ServerConfig {
	var serverConfig ServerConfig
	serverConfig.peers = make([]ServerPeer, len(clientConfigs))

	serverConfig.address = cliargs.Subnet.String()
	serverConfig.listenPort = cliargs.Port

	privateKey, publicKey, err := utils.NewX25519pair()
	if err != nil {
		panic(err)
	}

	serverConfig.privateKey = base64.StdEncoding.EncodeToString(privateKey)
	serverConfig.ServPublicKey = base64.StdEncoding.EncodeToString(publicKey)

	for i := range serverConfig.peers {
		serverConfig.peers[i].publicKey = clientConfigs[i].CliPublicKey
		serverConfig.peers[i].presharedKey = clientConfigs[i].PresharedKey
		serverConfig.peers[i].allowedIPs = clientConfigs[i].Address
	}

	return serverConfig
}
