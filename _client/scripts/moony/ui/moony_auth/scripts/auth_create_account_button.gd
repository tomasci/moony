extends Button


func _pressed() -> void:
	MoonyAuthState.moonyAuthIsAuthorized = false
	MoonyAuthState.moonyAuthShowCreateAccount = true
	return
