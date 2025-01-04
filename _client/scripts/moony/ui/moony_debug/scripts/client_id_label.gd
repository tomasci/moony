extends RichTextLabel


func _ready() -> void:
	self.visible = false
	return


func _process(delta: float) -> void:
	if MoonyAuthState.moonyAuthClientId:
		text = "[center]%s[/center]" % MoonyAuthState.moonyAuthClientId
		self.visible = true
	return
