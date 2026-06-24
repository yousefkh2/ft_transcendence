extends Control

@onready var role_label = $role

func _ready():
	role_label.text = "Connecting..."

	NetworkManager.connected.connect(_on_connected)
	NetworkManager.room_joined.connect(_on_room_joined)
	NetworkManager.connection_failed.connect(_on_connection_failed)

	NetworkManager.connect_to_server()


func _on_connected():
	print("Connected, joining room...")

	NetworkManager.join_room("TEST")


func _on_room_joined(role: String):
	print("Assigned role:", role)

	role_label.text = "Role: " + role


func _on_connection_failed(message: String):
	role_label.text = "Error: " + message
