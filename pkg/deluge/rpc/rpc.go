package rpc

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"errors"
	"io"
	"log"

	"github.com/gdm85/go-rencode"
)

const (
	rpcResponse = 1
	rpcError    = 2
	rpcEvent    = 3
)

// Connection represents a connetion to the server
type Connection struct {
	conn io.ReadWriteCloser
	open bool

	responses map[uint32]chan RPCResponseData
	errors    map[uint32]chan RPCErrorData

	events chan RPCEventData
}

type RPCResponseData interface{}

type RPCErrorData struct {
	ExceptionType    string
	ExceptionMessage string
	Traceback        string
}

type RPCEventData struct {
	EventName string
	Data      []interface{}
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
		conn:      conn,
		open:      true,
		responses: make(map[uint32]chan RPCResponseData),
		errors:    make(map[uint32]chan RPCErrorData),
		events:    make(chan RPCEventData),
	}

	go res.receive()

	return res, err
}

// Request sends a RPC Request to the server
func (c Connection) Request(id uint32, method string, args []interface{}, kwargs map[string]interface{}) (*RPCResponseData, *RPCErrorData) {
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

	c.responses[id] = make(chan RPCResponseData)
	c.errors[id] = make(chan RPCErrorData)

	c.conn.Write(message.Pack())

	var responseData *RPCResponseData
	var errrorData *RPCErrorData

	select {
	case res := <-c.responses[id]:
		responseData = &res
	case err := <-c.errors[id]:
		errrorData = &err
	}

	delete(c.responses, id)
	delete(c.errors, id)

	return responseData, errrorData
}

func (c *Connection) receive() {

	for c.open {
		m, _ := ReadMessage(c.conn)

		if m == nil {
			continue
		}

		compressed, _ := zlib.NewReader(bytes.NewBuffer(m.Body))
		decoder := rencode.NewDecoder(compressed)

		var body rencode.List
		decoder.Scan(&body)

		var messageType int
		body.Scan(&messageType)

		switch messageType {
		case rpcResponse:
			log.Println("Response!")

			var rid int
			var values RPCResponseData
			body.Scan(&messageType, &rid, &values)

			c.responses[uint32(rid)] <- values
		case rpcError:
			log.Println("Error!")

			var rid int
			var data RPCErrorData
			body.Scan(&messageType, &rid, &data.ExceptionType, &data.ExceptionMessage, &data.Traceback)

			c.errors[uint32(rid)] <- data
		case rpcEvent:
			log.Println("Event!")
		default:
			panic(errors.New("unknown message type"))
		}
	}
}
