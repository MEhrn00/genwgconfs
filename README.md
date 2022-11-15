# genwgconfs

Small cli tool to quickly generate basic Wireguard configs.

## Usage
```bash
genwgconfs -h
```

## Installation
#### Using `go install`
```bash
go install github.com/mehrn00/genwgconfs
```

Binary placed in `$GOPATH`.

#### Pre-built Releases
Pre-built binaries are in the [releases](https://github.com/MEhrn00/genwgconfs/releases) tab.

## Examples

Generates a server/client config pair and prints the result to stdout.
```bash
$ genwgconfs -endpoint vpn.example.com
# Server -------------------------------------------------------
[Interface]
Address = 10.0.0.0/24
PrivateKey = L2+h62C+eqPEGsJYlDzS8pdPEAQTm11u1g0hXOszmhM=
ListenPort = 58120

[Peer]
AllowedIPs = 10.0.0.1/32
PublicKey = WbKZUUhPJsUSyG7nTIRscQGoENqp4MDa8lLOCoXZNzE=
PresharedKey = kOzEvRw+WK9u8Er3Ng5kKYEyUCGUjMUew0sEZHWJVyc=

# Client 1 -----------------------------------------------------
[Interface]
Address = 10.0.0.1/24
PrivateKey = 8VhwStQvRZFvf8s4b8apXO0mtMvtNiDQutSfX56VcR4=

[Peer]
Endpoint = vpn.example.com:58120
PublicKey = 6kljBXqd93h1fc1qevG0feG/ldny7Ind7dlV/bIevhs=
PresharedKey = kOzEvRw+WK9u8Er3Ng5kKYEyUCGUjMUew0sEZHWJVyc=
AllowedIPs = 10.0.0.0/24
```

Generates VPN configs for 3 clients and outputs them into separate files with the `vpn` prefix.
```bash
$ genwgconfs -endpoint vpn.example.com -peers 3 -output vpn
$ ls
vpn_client1.conf  vpn_client2.conf  vpn_client3.conf  vpn_server.conf
```

Generates a VPN client/server pair but sets the `PersistentKeepalive` value to 25 and
forwards all traffic through the VPN by setting the `AllowedIPs` block to 0.0.0.0/0.
```bash
$ genwgconfs -endpoint vpn.example.com -pk 25 -all
# Server -------------------------------------------------------
[Interface]
Address = 10.0.0.0/24
PrivateKey = O3LVkoFlcAeLqH2pS+arqf5/R3mTjoJg/i1hjWjf2GU=
ListenPort = 58120

[Peer]
AllowedIPs = 10.0.0.1/32
PublicKey = VI66z9uJCuTnkjvmYiwiuyh7CoGYC42MjuFZE+YjZXQ=
PresharedKey = KQscjWjXKYrlF/HDjeFGnlFc45ZR9+ynqTfjtE5sM08=

# Client 1 -----------------------------------------------------
[Interface]
Address = 10.0.0.1/24
PrivateKey = 2EtmdCwcTqL9VWtFSnKBTedE3s98u3xg3st5FJZkiXo=

[Peer]
Endpoint = vpn.example.com:58120
PublicKey = hnyeFAVJTBOUJj8+1Wdk9HugED0rfb3YeDwZOqvXUwc=
PresharedKey = KQscjWjXKYrlF/HDjeFGnlFc45ZR9+ynqTfjtE5sM08=
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
```
