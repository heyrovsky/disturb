/*
Package p2p provides p2p functions for the application
*/

package p2p

import (
	"net"

	"github.com/heyrovsky/disturbdb/log"
	"github.com/heyrovsky/disturbdb/pkg/keys"
)

type Node struct {
	logger *log.Logger

	host net.IP
	port uint16

	privateKey keys.PrivateKey
	publicKey  keys.PublicKey
}
