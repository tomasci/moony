package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"moony/moony/core/event_dispatcher"
	"moony/moony/core/plugins"
	"moony/utils"
	"moony/utils/response"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
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

type MessageData struct {
	Plugin string      `json:"plugin"`
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

var dispatcher *event_dispatcher.EventDispatcher

// init is go predefined function, it executes before main
func init() {
	log.Println("Server init")
	// get dispatcher
	dispatcher = event_dispatcher.GetGlobalDispatcher()

	// get executable directory
	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatalf("failed to get executable directory: %v\n", err)
	}

	// try to load plugins
	pluginsDir := filepath.Join(exeDir, "plugins")
	if loadedPluginsCount, err := plugins.LoadPlugins(pluginsDir, dispatcher); err != nil {
		log.Fatalf("failed to load plugins: %v\n", err)
	} else {
		log.Printf("Loaded %d plugins", loadedPluginsCount)
	}
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

	ctx := context.Background()

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
		dispatcher.Dispatch("OnServerStarted", ctx, nil)
		// but you can remove this one :)
		dispatcher.Dispatch("CustomHelloWorldEvent", ctx, []any{"Hello world!"})

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
	dispatcher.Dispatch("hello_world_plugin.sum", ctx, []any{5, 7})
	dispatcher.Dispatch("hello_world_plugin.sum", ctx, []any{10, 15})
	dispatcher.Dispatch("hello_world_plugin.sum", ctx, []any{32461, 132})
	dispatcher.Dispatch("hello_world_plugin.sum", ctx, []any{5, "b"})
	dispatcher.RegisterEventHandler("hello_world_plugin.sum.result", func(ctx context.Context, data []any) {
		log.Println("Dispatcher: hello_world_plugin.sum.result (Example)", data)
	})
	dispatcher.Dispatch("hello_world_plugin.capitalize", ctx, []any{"welcome evening, sydney! i'm taylor!"})
	dispatcher.RegisterEventHandler("hello_world_plugin.capitalize.result", func(ctx context.Context, data []any) {
		log.Println("Dispatcher: hello_world_plugin.capitalize.result (Example)", data)
	})

	// create channel for os Signal values, can store only one signal
	sigChan := make(chan os.Signal, 1)
	// listen for SIGINT & SIGTERM and direct signals to sigChan
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// wait for system signal (interruption|termination signal)
	<-sigChan

	// dispatch server stopped event
	// don't delete this event because it may affect some code or plugins
	dispatcher.Dispatch("OnServerStopped", ctx, nil)

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

	var messageData MessageData
	err := json.Unmarshal(packet.data, &messageData)
	if err != nil {
		log.Println("failed to unmarshal packet:", err)

		responseJson, responseError := response.Error[any](500, "", "", nil, err)
		if responseError != nil {
			log.Println("failed to marshall response:", responseError)
		}

		_, udpError := conn.WriteToUDP(responseJson, packet.address)
		if udpError != nil {
			log.Println("failed to write response:", udpError)
		}
	}

	responseJson, responseError := response.Success[any]("", "", "Hello, Moony!")
	if responseError != nil {
		log.Println("failed to marshall response:", responseError)
	}

	_, udpError := conn.WriteToUDP(responseJson, packet.address)
	if udpError != nil {
		log.Println("failed to write response:", udpError)
	}

	// todo:
	// 1. get plugin by name
	// 2. call plugin method by name
}
