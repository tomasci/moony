extends Control

func _ready() -> void:
	self.visible = true
	
	var callable = Callable(self, "_onMoonyAuthLoginResult")
	MoonyClient.connect("moony_message_auth_login_result", callable)
	
	return


func _process(delta: float) -> void:
	if MoonyAuthState.moonyAuthIsAuthorized or MoonyAuthState.moonyAuthShowCreateAccount:
		self.visible = false
	else:
		self.visible = true
	
	return


func _onMoonyAuthLoginResult(data):
	var clientId = data[0]
	# todo: save user token
	MoonyAuthState.moonyAuthIsAuthorized = true
	MoonyAuthState.moonyAuthClientId = clientId
	return
