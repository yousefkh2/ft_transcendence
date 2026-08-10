extends Control

@onready var create_lobby_interface: Panel = $create_lobby_interface
@onready var lobby_interface: Panel = $lobby_interface

#func _ready():
	#role_label.text = "Connecting..."
#
	#NetworkManager.connected.connect(_on_connected)
	#NetworkManager.room_joined.connect(_on_room_joined)
	#NetworkManager.connection_failed.connect(_on_connection_failed)
#
	#NetworkManager.connect_to_server()
	
func _ready():
		create_lobby_interface.visible = false

func _on_back_button_pressed() -> void:
	create_lobby_interface.visible = false



func _on_create_lobby_pressed() -> void:
	create_lobby_interface.visible = true


func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/main_menu/menu.tscn")
