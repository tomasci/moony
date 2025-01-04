extends Button

@export var usernameInput: LineEdit
@export var passwordInput: LineEdit

func _pressed() -> void:
	var username = usernameInput.text
	var password = passwordInput.text 
	
	if username != "" and password != "":
		MoonyClient.sendMessage("auth", "login", [username, password])
	
	return
