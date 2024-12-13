extends Node

# this function is for debug purposes only
func _ready() -> void:
	print("moony state loaded")
	return

# in this state I will save all client and server related state 
# but it is not meant to save any of your data (messages and data between server and client)
# select another location in your project where you will save everything else
var mgcClientConnectionId: String
