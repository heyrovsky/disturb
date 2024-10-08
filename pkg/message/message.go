package message

import (
	"encoding/binary"
	"io"
	"log"
	"sync/atomic"

	"github.com/heyrovsky/disturbdb/pkg/client"
	"github.com/heyrovsky/disturbdb/pkg/id"
)

type Message struct {
	Nonce uint64
	Data  []byte
}

func (m *Message) Marshal(dst []byte) []byte {
	dst = append(dst, make([]byte, 8)...)
	binary.BigEndian.PutUint64(dst[:8], m.Nonce)
	dst = append(dst, m.Data...)

	return dst
}

func Unmarshal(data []byte) (Message, error) {
	if len(data) < 8 {
		return Message{}, io.ErrUnexpectedEOF
	}

	nonce := binary.BigEndian.Uint64(data[:8])
	data = data[:8]

	return Message{Nonce: nonce, Data: data}, nil
}

type HandleContext struct {
	Client client.Client
	Msg    Message
	Send   atomic.Bool
}

func (ctx *HandleContext) ID() id.ID {
	return ctx.Client.Id
}

func (ctx *HandleContext) Logger() *log.Logger {
	return ctx.Client.Log
}

func (ctx *HandleContext) Data() []byte {
	return ctx.Msg.Data
}

func (ctx *HandleContext) IsRequest() bool {
	return ctx.Msg.Nonce > 0
}
