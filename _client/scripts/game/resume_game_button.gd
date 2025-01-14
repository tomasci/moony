extends Button

func _ready() -> void:
	if !MoonyAuthState.moonyAuthIsAuthorized:
		self.visible = false
	return

func _process(delta: float) -> void:
	if MoonyAuthState.moonyAuthIsAuthorized and GameState.isPaused:
		self.visible = true
	else:
		self.visible = false
	return

func _pressed() -> void:
	Input.set_mouse_mode(Input.MOUSE_MODE_CAPTURED)
	GameState.isPaused = false
	return
