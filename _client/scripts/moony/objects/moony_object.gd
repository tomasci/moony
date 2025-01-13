extends Node3D

@export var label3d: Label3D

func _process(delta: float) -> void:
	var id = get_meta("id")
	if id:
		label3d.text = id
	pass
