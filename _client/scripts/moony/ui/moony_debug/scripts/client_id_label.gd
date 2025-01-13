extends RichTextLabel


func _ready() -> void:
	self.visible = false
	return


func _process(delta: float) -> void:
	if MoonyAuthState.moonyAuthClientId:
		text = "[right]%s[/right]" % MoonyAuthState.moonyAuthClientId
		self.visible = true
	return
