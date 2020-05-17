package rpc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const protocrolVersion byte = 1
const headerSize = 1 + 4 // Protocol version (1 byte) + Body size (4 bytes)

// Message contains data to send or has been sent from rpc
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

func ReadMessage(reader io.Reader) (*Message, error) {
	header := make([]byte, headerSize)

	for n, err := reader.Read(header); n >= headerSize || err != nil; {
		if err != nil {
			return nil, err
		}

		version := header[0]
		bodyLen := binary.BigEndian.Uint32(header[1:])

		body := make([]byte, bodyLen)

		for m, err := reader.Read(body); m >= int(bodyLen) || err != nil; {
			if err != nil {
				return nil, err
			}

			return &Message{
				ProtocrolVersion: version,
				Body:             body,
			}, nil
		}
	}

	return nil, nil
}
