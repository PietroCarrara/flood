package rpc

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"io"
	"log"

	"github.com/gdm85/go-rencode"
)

// Connection represents a connetion to the server
type Connection struct {
	conn io.ReadWriteCloser
	open bool
}

// Close closes the connection
func (c *Connection) Close() {
	c.open = false
	c.conn.Close()
}

// Connect opens a connection to communicate with the server
func Connect(address string) (*Connection, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})

	res := &Connection{
		conn: conn,
		open: true,
	}

	go res.read()

	return res, err
}

// Request sends a RPC Request to the server
func (c Connection) Request(id uint32, method string, args []interface{}, kwargs map[string]interface{}) {
	buf := &bytes.Buffer{}

	dict := rencode.Dictionary{}
	for key, value := range kwargs {
		dict.Add(key, value)
	}

	compressed, _ := zlib.NewWriterLevel(buf, zlib.NoCompression)
	encoder := rencode.NewEncoder(compressed)
	encoder.Encode(rencode.NewList(rencode.NewList(id, method, rencode.NewList(args...), dict)))
	compressed.Close()

	message := NewMessage(buf.Bytes())

	c.conn.Write(message.Pack())
}

func (c *Connection) read() {

	buffer := make([]byte, 1024)

	for c.open {
		n, _ := c.conn.Read(buffer)

		if n > 0 {
			log.Println(string(buffer[:n]))
		}
	}
}
