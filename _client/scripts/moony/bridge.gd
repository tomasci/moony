extends Node

# objects
var sphereObject = preload("res://scenes/objects/sphere.tscn")

# materials
var materials = {
	"red": preload("res://scenes/objects/materials/red.tres"),
	"green": preload("res://scenes/objects/materials/green.tres"),
	"blue": preload("res://scenes/objects/materials/blue.tres")
}

# object registry
var objectRegistry: Dictionary = {}

func _ready() -> void:
	print("moony bridge loaded")
	
	MoonyClient.connect("moony_message_hello_world_spawn_object_result", _onSpawnObject)
	return 

func _onSpawnObject(data):
	var object = data[0]
	var atPos = data[1]
	var material = data[2]
	
	print("_onSpawnObject", object, atPos, material)
	
	_spawnObject(object["id"], Vector3(atPos.x, atPos.y, atPos.z), material)
	return

func _spawnObject(objectId: String, position: Vector3, material: String):
	var instance = sphereObject.instantiate()
	var currentScene = get_tree().current_scene
	
	if instance is Node3D:
		objectRegistry[objectId] = instance
		instance.set_meta("id", objectId)
		instance.position = position
		
		var meshInstance = instance.get_node("Model")
	
		#if meshInstance and meshInstance is CSGSphere3D:
			#print("instance2", meshInstance.material)
		if materials[material]:
			meshInstance.material = materials[material]
		
	if currentScene:
		currentScene.add_child(instance)

	return
