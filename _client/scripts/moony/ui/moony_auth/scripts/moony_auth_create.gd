extends Control

@export var usernameInput: LineEdit
@export var passwordInput: LineEdit
@export var emailInput: LineEdit

func _ready() -> void:
	self.visible = false
	
	var callable = Callable(self, "_onMoonyAuthCreateResult")
	MoonyClient.connect("moony_message_auth_create_result", callable)
	
	return


func _process(delta: float) -> void:
	if !MoonyAuthState.moonyAuthIsAuthorized and MoonyAuthState.moonyAuthShowCreateAccount:
		self.visible = true
	else:
		self.visible = false
	
	return


func _onMoonyAuthCreateResult(data):
	# reset inputs
	usernameInput.text = ""
	passwordInput.text = ""
	emailInput.text = ""
		
	MoonyAuthState.moonyAuthIsAuthorized = false
	MoonyAuthState.moonyAuthShowCreateAccount = false
	return
