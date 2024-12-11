class_name MoonyClientNode
extends Node

var udp := PacketPeerUDP.new()

# udp callbacks 
# signals here
# udp callbacks end

func _ready() -> void:
	print("moony client loaded")
	udp.connect_to_host(MoonyConfig.mgcUDPAddress, MoonyConfig.mgcUDPPort)
	return

func _process(delta: float) -> void:
	if udp.is_socket_connected():
		if udp.get_available_packet_count() > 0:
			var packet = udp.get_packet()
			_onIncomingPacket(packet)
	return

func _onIncomingPacket(incomingPacket) -> void:
	print("onIncomingPacket: ", incomingPacket)
	return

func sendMessage(data) -> void:
	print("MoonyClient sendMessage data: ", data)
	
	var dataStr = JSON.stringify(data)
	var dataBase64 = Marshalls.utf8_to_base64(dataStr)
	var dataBin = Marshalls.base64_to_raw(dataBase64)
	print("MoonyClient sendMessage binary: ", dataBin)
	
	udp.put_packet(dataBin)
	print("MoonyClient sendMessage message sent")
	return
