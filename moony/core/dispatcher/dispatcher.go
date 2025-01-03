package dispatcher

// if you need "raw events" you can use dispatcher
// but if you don't - just use events wrapper

import (
	"context"
	"net"
	"sync"
)

// EventHandler - event handler type
type EventHandler func(ctx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any)

type EventDispatcher struct {
	// sync access to shared resources, thread-safe read-write access to handlers
	mu sync.RWMutex
	// map with eventType eventHandler association
	handlers map[string][]EventHandler
}

// pointer to singleton instance of EventDispatcher
var globalDispatcher *EventDispatcher

// once ensures that singleton globalDispatcher is created only once
var once sync.Once

// GetGlobalDispatcher provides a thread-safe way to get global instance of EventDispatcher
func GetGlobalDispatcher() *EventDispatcher {
	// once do used to initialize globalDispatcher on its first invocation
	// ensuring that the initialization code runs only once
	once.Do(func() {
		globalDispatcher = &EventDispatcher{
			handlers: make(map[string][]EventHandler),
		}
	})

	return globalDispatcher
}

// RegisterEventHandler - this is a method with receiver, basically method of EventDispatcher class if you're coming from other languages
// think about it like any class (like classes in TypeScript)
// so in TS it would be:
//
//	class EventDispatcher {
//		private handlers = {}
//	 public RegisterEventHandler(eventType: string, handler: ...) {
//			handlers[eventType] = handler
//	 }
//	 ... etc
//	}
//
// and now you can call it
// const dispatcher = new EventDispatcher()
// dispatcher.RegisterEventHandler(...)
// hope you understand it, because it took a while for me...
func (d *EventDispatcher) RegisterEventHandler(eventType string, handler EventHandler) {
	// lock dispatcher to add new handler
	d.mu.Lock()
	// unlock after updating
	defer d.mu.Unlock()
	// add new handler
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch - trigger all handlers by eventType
// same thing with receiver here and in any other cases
func (d *EventDispatcher) Dispatch(eventType string, ctx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
	// read lock to ensure safe concurrent reads from handlers map
	// basically, no new handlers will be added until current operation is done
	d.mu.RLock()
	// unlock after updating
	defer d.mu.RUnlock()

	// check if handlers exists for eventType
	if handlers, exist := d.handlers[eventType]; exist {
		// go through all handlers and call them
		for _, handler := range handlers {
			// calling each in new goroutine (allowing async event processing)
			go handler(ctx, conn, address, data)
		}
	}
}
