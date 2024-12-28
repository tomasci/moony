extends Control

# simple chat example 

# elements to work with
@export var chatMessagePrefab: Label
@export var chatMessageListContainer: Container

func _ready() -> void:
	# create callable callback for moony_message_any signal
	var callableOnMoonyMessageAny = Callable(self, "_onMoonyMessageAny")
	# connect signal with callback
	MoonyClient.connect("moony_message_hello_world_capitalize", callableOnMoonyMessageAny)
	return

# callback for moony_message_any signal
func _onMoonyMessageAny(data):
	print("_onMoonyMessageAny: ", data)
	# in this example just call append immediately 
	_appendChatMessage(data)
	return

func _appendChatMessage(data):
	# duplicate message prefab
	var chatMessageElement = chatMessagePrefab.duplicate()
	# set data as message text
	chatMessageElement.text = str(data)
	# and show element, because it is hidden by default (you don't want to see prefab on screen right?)
	chatMessageElement.show()
	
	# add message element to container 
	chatMessageListContainer.add_child(chatMessageElement)
	# and move it
	chatMessageListContainer.move_child(chatMessageElement, 0)
	return
