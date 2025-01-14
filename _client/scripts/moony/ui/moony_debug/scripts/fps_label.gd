extends RichTextLabel


func _ready() -> void:
	self.visible = false
	return


func _process(delta: float) -> void:
	var fps = Engine.get_frames_per_second()
	text = "[right]%s FPS[/right]" % fps
	self.visible = true
	return
