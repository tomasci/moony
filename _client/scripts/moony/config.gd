extends Node

# this function is for debug purposes only
func _ready() -> void:
	print("moony config loaded")
	return

# server connection details
# you can put your values directly here
# or use them from elsewhere
# this example shows how to use values from global state

# server ip address
const udpAddress: String = Globals.udpAddress
# server port 
const udpPort: int = Globals.udpPort
