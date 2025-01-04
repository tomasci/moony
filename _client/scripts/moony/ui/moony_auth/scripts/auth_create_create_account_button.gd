extends Button

@export var usernameInput: LineEdit
@export var passwordInput: LineEdit
@export var emailInput: LineEdit

func _pressed() -> void:
	var username = usernameInput.text
	var password = passwordInput.text
	var email = emailInput.text
	
	if username != "" and password != "" and email != "":
		MoonyClient.sendMessage("auth", "create", [username, password, email])
	
	return
