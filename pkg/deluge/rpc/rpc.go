package rpc

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"errors"
	"fmt"
	"io"

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

type RPCResponseData struct {
	rencode.List
}

type RPCErrorData struct {
	ExceptionType    string
	ExceptionMessage string
	ExceptionKwargs  map[string]interface{}
	Traceback        string
}

type RPCEventData struct {
	EventName string
	Data      []interface{}
}

func (err *RPCErrorData) Error() string {
	if err.ExceptionMessage != "" {
		return fmt.Sprintf("%s: %v", err.ExceptionType, err.ExceptionMessage)
	}

	return err.ExceptionType
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

// Request sends a RPC request that has no kwargs
func (c Connection) Request(id uint32, method string, args ...interface{}) (*RPCResponseData, *RPCErrorData) {
	return c.RequestArgsKwargs(id, method, args, map[string]interface{}{})
}

// Request sends a RPC request that has kwargs
func (c Connection) RequestKwargs(id uint32, method string, kwargs map[string]interface{}, args ...interface{}) (*RPCResponseData, *RPCErrorData) {
	return c.RequestArgsKwargs(id, method, args, kwargs)
}

// RequestArgsKwargs sends a RPC Request to the server using args and kwargs
func (c Connection) RequestArgsKwargs(id uint32, method string, args []interface{}, kwargs map[string]interface{}) (*RPCResponseData, *RPCErrorData) {
	buf := &bytes.Buffer{}

	dict := rencode.Dictionary{}
	for key, value := range kwargs {
		dict.Add(key, value)
	}

	// TODO: Fix go-rencode aversion to maps
	for i, val := range args {
		switch v := val.(type) {
		case map[string]interface{}:
			dict := rencode.Dictionary{}
			for key, value := range v {
				dict.Add(key, value)
			}
			args[i] = dict
		}
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
			var rid int
			var values interface{}
			body.Scan(&messageType, &rid, &values)

			c.responses[uint32(rid)] <- RPCResponseData{rencode.NewList(values)}
		case rpcError:
			var rid int
			var data RPCErrorData
			var args interface{}
			var kwargs rencode.Dictionary
			body.Scan(&messageType, &rid, &data.ExceptionType, &args, &kwargs, &data.Traceback)

			switch v := args.(type) {
			case rencode.List:
				v.Scan(&data.ExceptionMessage, &data.ExceptionType, &data.Traceback)
			case string:
				data.ExceptionMessage = v
			default:
				panic("exception message of type unknown")
			}

			c.errors[uint32(rid)] <- data
		case rpcEvent:
			// TODO
		default:
			panic(errors.New("unknown message type"))
		}
	}
}
