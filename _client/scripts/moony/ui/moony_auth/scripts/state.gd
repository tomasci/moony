extends Node

func _ready() -> void:
	print("moony auth state loaded")
	return

var moonyAuthIsAuthorized: bool = false
var moonyAuthToken: String
var moonyAuthClientId: String
var moonyAuthShowCreateAccount: bool = false
