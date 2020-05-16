package main

import (
	"time"

	"github.com/PietroCarrara/flood/pkg/deluge/rpc"
)

// Flood keeps state about the deluge server
// and is able to interact with it
type Flood struct {
}

func main() {
	conn, _ := rpc.Connect("localhost:58846")
	defer conn.Close()

	conn.Request(1, "daemon.info", []interface{}{}, map[string]interface{}{})

	time.Sleep(2 * time.Second)
}
