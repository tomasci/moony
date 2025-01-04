extends Button


func _pressed() -> void:
	MoonyAuthState.moonyAuthIsAuthorized = false 
	MoonyAuthState.moonyAuthShowCreateAccount = false
	return
