package client

import (
	"container/list"
	"sync"
)

type clientMapEntry struct {
	el     *list.Element
	client *Client
}

type ClientMap struct {
	sync.Mutex

	cap     uint
	order   *list.List
	entries map[string]clientMapEntry
}

func newClientMap(cap uint) *ClientMap {
	return &ClientMap{
		cap:     cap,
		order:   list.New(),
		entries: make(map[string]clientMapEntry, cap),
	}
}
