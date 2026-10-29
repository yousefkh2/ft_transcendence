extends Control

const MAIN_LOBBY_SCENE := "res://games/lobby/lobby.tscn"
const STARTING_SECONDS := 15
const GAME_WEBSOCKET_URL := "ws://localhost:8080/ws"

@onready var game_timer: Timer = $GameTimer
@onready var timer_label: Label = $TimerLabel
@onready var leave_button: Button = $LeaveButton

var game_socket := WebSocketPeer.new()
var remaining_seconds := STARTING_SECONDS
var defeat_popup: PopupPanel
var returning_to_lobby := false
var game_room_join_sent := false

func _ready() -> void:
	game_timer.timeout.connect(_on_game_timer_timeout)
	_update_timer_label()
	game_timer.start()
	var error := game_socket.connect_to_url(GAME_WEBSOCKET_URL)
	if error != OK:
		push_error("Failed to connect to the game server: %s" % error)

func _on_game_timer_timeout() -> void:
	remaining_seconds -= 1
	_update_timer_label()
	if remaining_seconds <= 0:
		_show_defeat_popup()

func _on_leave_button_pressed() -> void:
	_return_to_main_lobby()

func _process(_delta: float) -> void:
	game_socket.poll()
	if game_socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		return
	if not game_room_join_sent:
		_join_game_room()
		game_room_join_sent = true

	while game_socket.get_available_packet_count() > 0:
		var response = JSON.parse_string(game_socket.get_packet().get_string_from_utf8())
		if typeof(response) != TYPE_DICTIONARY:
			continue
		if response.get("type", "") == "room.joined":
			continue
		if response.get("type", "") == "player.disconnected":
			GameState.lobby_message = str(response.get("message", "A player left the game"))
			_return_to_main_lobby()

func _join_game_room() -> void:
	game_socket.send_text(JSON.stringify({
		"type": "room.join",
		"roomCode": GameState.lobby_code,
		"token": GameState.auth_token,
	}))

func _update_timer_label() -> void:
	timer_label.text = str(remaining_seconds)

func _show_defeat_popup() -> void:
	game_timer.stop()
	leave_button.disabled = true

	defeat_popup = PopupPanel.new()
	defeat_popup.size = Vector2(420, 220)
	defeat_popup.exclusive = true

	var contents := VBoxContainer.new()
	contents.add_theme_constant_override("separation", 24)
	contents.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT, Control.PRESET_MODE_MINSIZE, 24)

	var title := Label.new()
	title.text = "Defeated"
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 48)
	contents.add_child(title)

	var exit_button := Button.new()
	exit_button.text = "Exit to lobby"
	exit_button.custom_minimum_size = Vector2(0, 56)
	exit_button.pressed.connect(_return_to_main_lobby)
	contents.add_child(exit_button)

	defeat_popup.add_child(contents)
	add_child(defeat_popup)
	defeat_popup.popup_centered()

func _return_to_main_lobby() -> void:
	if returning_to_lobby:
		return
	returning_to_lobby = true
	game_timer.stop()
	if game_socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		game_socket.close()
	get_tree().change_scene_to_file(MAIN_LOBBY_SCENE)
