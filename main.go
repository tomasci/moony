package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"moony/moony/core"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	localAddress = "127.0.0.1"
	hostAddress  = "0.0.0.0"
	port         = 25505
)

type PacketData struct {
	address *net.UDPAddr
	data    []byte
}

var dispatcher *core.EventDispatcher

func init() {
	log.Println("Server init")
	dispatcher = core.GetGlobalDispatcher()
}

func main() {
	log.Println("Server is in startup")

	// host flag definition
	isHostPointer := flag.Bool("host", false, "Set this flag to listen on all interfaces (will start server at 0.0.0.0)")

	// parse flags
	flag.Parse()

	isHost := *isHostPointer
	var address string

	// check host flag and change address according to it
	if isHost {
		address = fmt.Sprintf("%s:%d", hostAddress, port)
	} else {
		address = fmt.Sprintf("%s:%d", localAddress, port)
	}

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

	// waitGroup to sync goroutines
	var wg sync.WaitGroup
	quit := make(chan struct{})
	// channel for packets; main loop send packets to channel to be processed by worker goroutines
	packetChan := make(chan PacketData, 100)

	// worker goroutines setup
	// starts a number of goroutines equal to threadsCount; each goroutine is a worker that processes packets
	// these workers started once, when server started, and alive until server stopped
	// no more workers created during execution
	for i := 0; i < threadsCount; i++ {
		wg.Add(1)
		// next: read worker func
		go worker(i, packetChan, &wg, quit, conn)
	}

	// main loop
	go func() {
		log.Println("Server started")
		// dispatch server started event
		// don't delete this event because it may affect some code or plugins
		dispatcher.Dispatch(core.OnServerStarted, context.Background(), nil)
		// but you can remove this one :)
		dispatcher.Dispatch("CustomHelloWorldEvent", context.Background(), "Hello, world!")

		for {
			// buffer for incoming data
			buffer := make([]byte, 1024)

			// reading data from UDP connection into a buffer
			n, clientAddress, err := conn.ReadFromUDP(buffer) // n - number of bytes in buffer
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Println("Server stopped: Connection closed")
					break
				}

				select {
				case <-quit:
					return // quit signaled
				default:
					log.Println("failed to read from UDP:", err)
				}
			}

			select {
			case packetChan <- PacketData{clientAddress, buffer[:n]}: // send packet to workers
				// each received packet is sent into packetChan
				// where they become available for processing by goroutines (and each goroutine is a separate worker?)
				// todo: remove, debug purposes only
				log.Printf("Incoming packet from %v (%d bytes): %v\n", clientAddress, n, string(buffer[:n]))
			case <-quit:
				return // break loop if quit signaled
			}
		}
	}()

	// these handlers are placed here just as an example
	// you can remove them (all three)
	dispatcher.RegisterEventHandler(core.OnServerStarted, func(ctx context.Context, data interface{}) {
		log.Println("Dispatcher: Server started (Example)", data)
	})

	dispatcher.RegisterEventHandler(core.OnServerStopped, func(ctx context.Context, data interface{}) {
		log.Println("Dispatcher: Server stopped (Example)", data)
	})

	dispatcher.RegisterEventHandler("CustomHelloWorldEvent", func(ctx context.Context, data interface{}) {
		log.Println("Dispatcher: CustomHelloWorldEvent (Example)", data)
	})

	// create channel for os Signal values, can store only one signal
	sigChan := make(chan os.Signal, 1)
	// listen for SIGINT & SIGTERM and direct signals to sigChan
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// wait for system signal (interruption|termination signal)
	<-sigChan

	// dispatch server stopped event
	// don't delete this event because it may affect some code or plugins
	dispatcher.Dispatch(core.OnServerStopped, context.Background(), nil)

	close(quit) // signal all goroutines to exit
	close(packetChan)
	wg.Wait() // wait for all workers to finish
	log.Println("Shutting down")
}

func worker(id int, packetChan chan PacketData, wg *sync.WaitGroup, quit <-chan struct{}, conn *net.UDPConn) {
	defer wg.Done() // signal to waitGroup that worker is done when function exits

	// continuously receive packets from packetChan and process them
	for {
		select {
		case packet := <-packetChan:
			processPacket(id, packet, conn)
		case <-quit:
			log.Printf("Worker %d stopped\n", id)
			return
		}
	}
}

func processPacket(id int, packet PacketData, conn *net.UDPConn) {
	log.Printf("Worker %d processing packet from %v: %v\n", id, packet.address, string(packet.data))

	// send same incoming packet back to client (just for demo)
	timestamp := time.Now().UnixMilli()
	response := []byte(fmt.Sprintf("%d", timestamp))
	_, err := conn.WriteToUDP(response, packet.address)
	if err != nil {
		log.Println("failed to write to UDP:", err)
	}
}
