extends Control

#func _ready():
	#role_label.text = "Connecting..."
#
	#NetworkManager.connected.connect(_on_connected)
	#NetworkManager.room_joined.connect(_on_room_joined)
	#NetworkManager.connection_failed.connect(_on_connection_failed)
#
	#NetworkManager.connect_to_server()


func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/main_menu/menu.tscn")
