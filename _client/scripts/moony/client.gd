class_name MoonyClientNode
extends Node

var udp := PacketPeerUDP.new()

# udp callbacks 
# register your signals here
# all signals must start with "moony_message_" prefix
# you need this to understand what is going on in your code and where server code is
signal moony_message_msync_ping_result
signal moony_message_hello_world_capitalize_result
signal moony_message_hello_world_spawn_object_result
signal moony_message_auth_login_result
signal moony_message_auth_create_result
# udp callbacks end

func _ready() -> void:
	print("moony client loaded")
	# "connect" to server (will not send any packets actually, makes is_socket_connected = true)
	udp.connect_to_host(MoonyConfig.udpAddress, MoonyConfig.udpPort)
	# and this command will send connect event to server
	MoonyClient.sendMessage("msync", "connect", [])
	return

func _notification(what) -> void:
	if what == NOTIFICATION_WM_CLOSE_REQUEST:
		print("Game closed.")
		# update server when user closes the game
		# it will not help, when user force quit game, but at least in this case
		MoonyClient.sendMessage("msync", "disconnect", [])
	return

func _process(delta: float) -> void:
	if udp.is_socket_connected():
		if udp.get_available_packet_count() > 0:
			# get packets
			var packet = udp.get_packet()
			# and process 
			_onIncomingPacket(packet)
	return

# process incoming packets here
func _onIncomingPacket(incomingPacket: PackedByteArray) -> void:
	print("onIncomingPacket: ", incomingPacket)
	
	# convert hex to base64
	var packetBase64 = Marshalls.raw_to_base64(incomingPacket)
	# convert base64 to string
	var packetString = Marshalls.base64_to_utf8(packetBase64)
	# parse as json object
	var packetObject = JSON.parse_string(packetString)
	print("_onIncomingPacket parsed object: ", packetBase64, packetString, packetObject)
	
	var status = packetObject["status"]
	
	if status == 200:
		var plugin = packetObject["plugin"]
		var method = packetObject["method"]
		var data = packetObject["data"]
		var signalName = "moony_message_" + str(plugin) + "_" + str(method)
	
		print("plugin: ", plugin)
		print("method: ", method)
		print("data: ", data)
		print("signalName: ", signalName)
	
		emit_signal(signalName, data)
	elif status == 500:
		print("error happened...")
		return
	else:
		print("unknown status")
		return
	
	return

# send messages to server
func sendMessage(plugin: String, method: String, data: Array) -> void:
	print("MoonyClient sendMessage data: ", data)
	
	var preparedMessage = {
		"plugin": plugin,
		"method": method,
		"data": data
	}
	
	# convert data to string 
	var dataStr = JSON.stringify(preparedMessage)
	# convert data to base64
	var dataBase64 = Marshalls.utf8_to_base64(dataStr)
	# convert data to binary
	var dataBin = Marshalls.base64_to_raw(dataBase64)
	print("MoonyClient sendMessage binary: ", dataBin)
	
	# send packet to server
	udp.put_packet(dataBin)
	print("MoonyClient sendMessage message sent")
	return
