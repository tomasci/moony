package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const (
	address = "0.0.0.0:25505"
)

func main() {
	log.Println("Server is in startup")

	// get logical CPU core count (threads count)
	threadsCount := runtime.NumCPU()
	// set number of threads that can run simultaneously
	runtime.GOMAXPROCS(threadsCount)

	log.Println("Threads:", threadsCount)

	// trying to resolve ip and port (basically parsing ip and port from string)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalln("failed to resolve address:", err)
	}

	log.Println("Resolved address:", addr)

	// creating udp socket connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln("failed to listen (failed to bind address):", err)
	}

	// cleanup (close connection when main function exits)
	defer func(conn *net.UDPConn) {
		log.Println("Connection closed")
		err := conn.Close()
		if err != nil {
			log.Fatalln("failed to close UDP connection:", err)
		}
	}(conn)

	log.Println("Listening on:", address)

	// create channel for os Signal values, can store only one signal
	sigChan := make(chan os.Signal, 1)
	// listen for SIGINT & SIGTERM and direct signals to sigChan
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// main loop
	go func() {
		log.Println("Server started")

		for {
			// buffer for incoming data
			buffer := make([]byte, 1024)

			// reading data from UDP connection into a buffer
			n, _, err := conn.ReadFromUDP(buffer) // n - number of bytes in buffer
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Println("Server stopped: Connection closed")
					break
				}

				log.Println("failed to read from UDP:", err)
				continue
			}

			// processing data (creating slice of data size)
			packet := make([]byte, n)
			// copy buffer (up to number of bytes) into packet
			copy(packet, buffer[:n])

			// todo: remove, debug purposes only
			log.Println("Incoming packet:", string(packet))
		}
	}()

	// wait for signal
	<-sigChan
	log.Println("Shutting down")
}
