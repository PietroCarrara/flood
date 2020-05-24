package rpc

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/PietroCarrara/rencode"
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

	responses map[int]chan ResponseData
	errors    map[int]chan ErrorData

	events chan EventData
}

type ResponseData []interface{}

type ErrorData struct {
	ExceptionType    string
	ExceptionMessage string
	ExceptionKwargs  map[string]interface{}
	Traceback        string
}

type EventData struct {
	EventName string
	Data      []interface{}
}

func (err *ErrorData) Error() string {
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
		responses: make(map[int]chan ResponseData),
		errors:    make(map[int]chan ErrorData),
		events:    make(chan EventData),
	}

	go res.receive()

	return res, err
}

// Request sends a RPC request that has no kwargs
func (c Connection) Request(id int, method string, args ...interface{}) (*ResponseData, error) {
	return c.RequestArgsKwargs(id, method, args, map[string]interface{}{})
}

// Request sends a RPC request that has kwargs
func (c Connection) RequestKwargs(id int, method string, kwargs map[string]interface{}, args ...interface{}) (*ResponseData, error) {
	return c.RequestArgsKwargs(id, method, args, kwargs)
}

// RequestArgsKwargs sends a RPC Request to the server using args and kwargs
func (c Connection) RequestArgsKwargs(id int, method string, args []interface{}, kwargs map[string]interface{}) (*ResponseData, error) {
	buf := &bytes.Buffer{}

	compressed, _ := zlib.NewWriterLevel(buf, zlib.NoCompression)

	call := rencode.List{
		rencode.List{
			id,
			method,
			args,
			kwargs,
		},
	}

	bytes, err := rencode.Encode(call)
	if err != nil {
		return nil, err
	}

	_, err = compressed.Write(bytes)
	if err != nil {
		return nil, err
	}
	compressed.Close()

	message := NewMessage(buf.Bytes())

	c.responses[id] = make(chan ResponseData)
	c.errors[id] = make(chan ErrorData)

	c.conn.Write(message.Pack())

	var responseData *ResponseData
	var errrorData error

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
		bytes, _ := ioutil.ReadAll(compressed)

		var body rencode.List
		rencode.Decode(bytes, &body)

		var messageType int
		rencode.ScanSlice(body, &messageType)

		switch messageType {
		case rpcResponse:
			var rid int
			var values interface{}
			rencode.ScanSlice(body, nil, &rid, &values)

			c.responses[rid] <- ResponseData{values}
		case rpcError:
			var rid int
			var data ErrorData
			var args interface{}
			var kwargs map[string]interface{}
			rencode.ScanSlice(body, nil, &rid, &data.ExceptionType, &args, &kwargs, &data.Traceback)

			if data.ExceptionType == "WrappedException" {
				rencode.ScanSlice(args, &data.ExceptionMessage, &data.ExceptionType, &data.Traceback)
			}

			c.errors[rid] <- data
		case rpcEvent:
			// TODO
		default:
			panic(errors.New("unknown message type"))
		}
	}
}
