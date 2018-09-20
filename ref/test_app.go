package main

import (
	"net"
	"sync"
)

var (
	// Conn holds the connection to proxy router socket object
	Conn net.Conn

	// ConnMutex used to lock the Conn object when used by a thread
	ConnMutex sync.Mutex
)

func init() {
	go SetupConnection()
}

// func main() {
// 	LocalProxy()
// }

// LocalProxy method get's new HTTP request and proxies it through the
// proxy routers
func LocalProxy() {
	if listen, lErr := net.Listen("tcp", "127.0.0.1:4121"); lErr == nil {
		for conn, cErr := listen.Accept(); cErr == nil; conn, cErr = listen.Accept() {
			go func(conn net.Conn) {
				defer conn.Close()

				reqBuff := make([]byte, 16384)
				resBuff := make([]byte, 16384)

				n, err := conn.Read(reqBuff)
				if err != nil {
					conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
					return
				}

				ConnMutex.Lock()
				Conn.Write(reqBuff[:n])
				n, err = Conn.Read(resBuff)
				ConnMutex.Unlock()

				if err != nil {
					conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
					return
				}

				conn.Write(resBuff[:n])
			}(conn)
		}
	}
}

// SetupConnection sets up the connection to the proxy router
func SetupConnection() {
	if conn, cErr := net.Dial("tcp", "127.0.0.1:4123"); cErr == nil {
		Conn = conn
	} else {
		panic("Proxy Router Unreachable")
	}
}
