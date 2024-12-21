# Server Events

Currently, there are only 2 events dispatched by the server:

* OnServerStarted - occurs when you start the server (listening started)
* OnServerStopped - occurs when you try to stop the server or the server is stopped by some signal.

No data is sent in these events.