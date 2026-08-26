extends Control

@onready var create_lobby_interface: Panel = $create_lobby_interface
@onready var lobby_interface: Panel = $lobby_interface
@onready var http_request: HTTPRequest = $HTTPRequest
@onready var lobby_code_input: LineEdit = $lobby_interface/LineEdit
@onready var message: Label = $message_interface/message
@onready var debug_menu: Panel = $debug_menu
@onready var username: Label = $debug_menu/username
@onready var password: Label = $debug_menu/password
@onready var email: Label = $debug_menu/email

func _ready():
	create_lobby_interface.visible = false

	http_request.login_success.connect(_on_login_success)
	http_request.login_failed.connect(_on_login_failed)
	http_request.lobby_created.connect(_on_lobby_created)
	http_request.lobby_creation_failed.connect(_on_lobby_creation_failed)
	http_request.lobby_joined.connect(_on_lobby_joined)
	http_request.lobby_join_failed.connect(_on_lobby_join_failed)

	http_request.login()
	
	

func _on_login_success() -> void:
	print("Logged in successfully.")
	username.text = GameState.player_name
	password.text = GameState.player_pass
	email.text = GameState.player_email

func _on_login_failed(message: String) -> void:
	push_error("Login failed: " + message)

func _on_back_button_pressed() -> void:
	create_lobby_interface.visible = false

func _on_create_lobby_pressed() -> void:
	create_lobby_interface.visible = true

func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/main_menu/menu.tscn")

func _on_select_dev_game_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/dev/dev_game_select.tscn")

func _on_create_pressed() -> void:
	http_request.create_lobby()

func _on_lobby_created(code: String) -> void:
	print("Lobby created with code: ", code)
	get_tree().change_scene_to_file("res://games/pregame_lobby/pregame_lobby.tscn")

func _on_lobby_creation_failed(message: String) -> void:
	push_error("Could not create lobby: " + message)

func _on_join_pressed() -> void:
	var code = lobby_code_input.text.strip_edges()
	if code == "":
		message.text = "Please enter a lobby code."
		return
	http_request.join_lobby(code)

func _on_lobby_joined(lobby_data: Dictionary) -> void:
	print("Joined lobby: ", lobby_data)
	get_tree().change_scene_to_file("res://games/pregame_lobby/pregame_lobby.tscn")

func _on_lobby_join_failed(message: String) -> void:
	push_error("Could not join lobby: " + message)
