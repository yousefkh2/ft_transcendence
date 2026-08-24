extends Control

@onready var create_lobby_interface: Panel = $create_lobby_interface
@onready var lobby_interface: Panel = $lobby_interface
@onready var http_request: HTTPRequest = $HTTPRequest  # matches the node name in your scene tree

func _ready():
	create_lobby_interface.visible = false
	http_request.lobby_created.connect(_on_lobby_created)
	http_request.lobby_creation_failed.connect(_on_lobby_creation_failed)

func _on_back_button_pressed() -> void:
	create_lobby_interface.visible = false

func _on_create_lobby_pressed() -> void:
	create_lobby_interface.visible = true

func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/main_menu/menu.tscn")

func _on_select_dev_game_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/dev/dev_game_select.tscn")

# This is your "create" button inside create_lobby_interface
func _on_create_pressed() -> void:
	http_request.create_lobby()

func _on_lobby_created(code: String) -> void:
	print("Lobby created with code: ", code)
	# stash the code somewhere the next scene can read it, e.g. an autoload:
	# GameState.current_lobby_code = code
	get_tree().change_scene_to_file("res://games/pregame_lobby/pregame_lobby.tscn")

func _on_lobby_creation_failed(message: String) -> void:
	push_error("Could not create lobby: " + message)
	# maybe show an error label in create_lobby_interface here
