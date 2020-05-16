package rpc

import (
	"bytes"
	"encoding/binary"
)

const protocrolVersion = 1

// Message contains data to be sent to rpc
type Message struct {
	ProtocrolVersion byte
	Body             []byte
}

// NewMessage creates a new message
func NewMessage(body []byte) *Message {
	return &Message{
		ProtocrolVersion: protocrolVersion,
		Body:             body,
	}
}

// Pack packs a message to be sent to rpc
func (m Message) Pack() []byte {
	buffer := &bytes.Buffer{}

	binary.Write(buffer, binary.BigEndian, m.ProtocrolVersion)
	binary.Write(buffer, binary.BigEndian, uint32(len(m.Body)))
	binary.Write(buffer, binary.BigEndian, m.Body)

	return buffer.Bytes()
}
