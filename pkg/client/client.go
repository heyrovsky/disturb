package client

import (
	"log"

	"github.com/heyrovsky/disturbdb/p2p"
	"github.com/heyrovsky/disturbdb/pkg/id"
)

type Client struct {
	Node *p2p.Node

	Id id.ID

	Log *log.Logger
}
