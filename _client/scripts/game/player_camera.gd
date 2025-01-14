extends Camera3D

@export var mouse_sensitivity: float = 0.15
@export var move_speed: float = 5.0
@export var vertical_speed: float = 5.0  # Add a separate speed for vertical movement

var velocity: Vector3 = Vector3()
var input_rotation: Vector3 = Vector3()

func _input(event):
	if event is InputEventKey and event.pressed:
		if event.keycode == KEY_ESCAPE:
			if GameState.isPaused:
				GameState.isPaused = false
				Input.set_mouse_mode(Input.MOUSE_MODE_CAPTURED)
			else:
				GameState.isPaused = true
				Input.set_mouse_mode(Input.MOUSE_MODE_VISIBLE)
	
	if GameState.isPaused:
		return

	if event is InputEventMouseMotion:
		input_rotation.y -= event.relative.x * mouse_sensitivity
		input_rotation.x -= event.relative.y * mouse_sensitivity
		input_rotation.x = clamp(input_rotation.x, -90, 90)
		rotation_degrees = Vector3(input_rotation.x, input_rotation.y, 0)

	if event is InputEventMouseButton:
		if Input.get_mouse_mode() == Input.MOUSE_MODE_VISIBLE:
			Input.set_mouse_mode(Input.MOUSE_MODE_CAPTURED)

func _process(delta):
	if GameState.isPaused:
		return

	velocity = Vector3()
	if Input.is_action_pressed("move_forward"):
		velocity.z -= 1
	if Input.is_action_pressed("move_backward"):
		velocity.z += 1
	if Input.is_action_pressed("move_left"):
		velocity.x -= 1
	if Input.is_action_pressed("move_right"):
		velocity.x += 1
		
	#velocity = velocity.normalized() * move_speed
	#velocity = global_transform.basis * velocity
	#global_translate(velocity * delta)
	
	# Apply horizontal movement relative to the player's orientation
	velocity = global_transform.basis * velocity.normalized() * move_speed * delta

	# Vertical movement (absolute in world space)
	var vertical_velocity = Vector3()
	if Input.is_action_pressed("move_up"):
		vertical_velocity.y += 1
	if Input.is_action_pressed("move_down"):
		vertical_velocity.y -= 1

	# Apply vertical movement
	vertical_velocity = vertical_velocity.normalized() * vertical_speed * delta

	# Combine horizontal and vertical movement
	global_translate(velocity + vertical_velocity)
