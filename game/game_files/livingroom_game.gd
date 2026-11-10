extends Control

const MAIN_LOBBY_SCENE := "res://games/lobby/lobby.tscn"
const STARTING_SECONDS := 15
const GAME_WEBSOCKET_URL := "ws://localhost:8080/ws"
const GAME_FONT: FontFile = preload("res://game_files/fonts/GrapeSoda.ttf")
const ACCENT_RED := Color(0.7529412, 0.22745098, 0.2, 1)
const ACCENT_YELLOW := Color(1, 0.76862746, 0.20392157, 1)

@onready var game_timer: Timer = $GameTimer
@onready var timer_label: Label = $TimerLabel
@onready var leave_button: Button = $LeaveButton

var game_socket := WebSocketPeer.new()
var remaining_seconds := STARTING_SECONDS
var defeat_popup: Control
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
			# Once time is up the round is over: let each player leave the
			# Time's up screen on their own instead of pulling both out.
			if defeat_popup:
				continue
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

	# Full-screen dimmer that also blocks clicks/drags on the game underneath.
	defeat_popup = Control.new()
	defeat_popup.set_anchors_preset(Control.PRESET_FULL_RECT)
	defeat_popup.mouse_filter = Control.MOUSE_FILTER_STOP

	var dimmer := ColorRect.new()
	dimmer.color = Color(0, 0, 0, 0.65)
	dimmer.set_anchors_preset(Control.PRESET_FULL_RECT)
	dimmer.mouse_filter = Control.MOUSE_FILTER_IGNORE
	defeat_popup.add_child(dimmer)

	var center := CenterContainer.new()
	center.set_anchors_preset(Control.PRESET_FULL_RECT)
	center.mouse_filter = Control.MOUSE_FILTER_IGNORE
	defeat_popup.add_child(center)

	var card := PanelContainer.new()
	card.custom_minimum_size = Vector2(760, 0)
	card.add_theme_stylebox_override("panel", _make_box(Color(1, 1, 1, 1), 24, 12, 48))
	center.add_child(card)

	var contents := VBoxContainer.new()
	contents.alignment = BoxContainer.ALIGNMENT_CENTER
	contents.add_theme_constant_override("separation", 28)
	card.add_child(contents)

	var title := Label.new()
	title.text = tr("_GAME_TIMES_UP_")
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_override("font", GAME_FONT)
	title.add_theme_font_size_override("font_size", 110)
	title.add_theme_color_override("font_color", ACCENT_RED)
	title.add_theme_color_override("font_shadow_color", Color(0.07058824, 0.07058824, 0.07058824, 1))
	title.add_theme_constant_override("shadow_offset_x", 4)
	title.add_theme_constant_override("shadow_offset_y", 4)
	contents.add_child(title)

	var subtitle := Label.new()
	subtitle.text = tr("_GAME_TIMES_UP_TEXT_")
	subtitle.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	subtitle.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	subtitle.add_theme_font_override("font", GAME_FONT)
	subtitle.add_theme_font_size_override("font_size", 44)
	subtitle.add_theme_color_override("font_color", Color(0.15, 0.15, 0.15, 1))
	contents.add_child(subtitle)

	var exit_button := Button.new()
	exit_button.text = tr("_GAME_BACK_TO_LOBBY_")
	exit_button.custom_minimum_size = Vector2(420, 80)
	exit_button.size_flags_horizontal = Control.SIZE_SHRINK_CENTER
	exit_button.add_theme_font_override("font", GAME_FONT)
	exit_button.add_theme_font_size_override("font_size", 48)
	exit_button.add_theme_color_override("font_color", Color(0, 0, 0, 1))
	exit_button.add_theme_color_override("font_hover_color", Color(0, 0, 0, 1))
	exit_button.add_theme_color_override("font_pressed_color", Color(0, 0, 0, 1))
	exit_button.add_theme_stylebox_override("normal", _make_box(ACCENT_YELLOW, 16, 5, 0))
	exit_button.add_theme_stylebox_override("hover", _make_box(Color(0.9079417, 0.9079417, 0.90794164, 1), 16, 5, 0))
	exit_button.add_theme_stylebox_override("pressed", _make_box(Color(0.79836667, 0.79836625, 0.7983663, 1), 16, 5, 0))
	exit_button.add_theme_stylebox_override("focus", StyleBoxEmpty.new())
	exit_button.pressed.connect(_return_to_main_lobby)
	contents.add_child(exit_button)

	add_child(defeat_popup)

	# Fade the dimmer in and pop the card up.
	dimmer.modulate.a = 0.0
	card.resized.connect(func(): card.pivot_offset = card.size / 2.0)
	card.scale = Vector2(0.6, 0.6)
	card.modulate.a = 0.0
	var tween := create_tween().set_parallel(true)
	tween.tween_property(dimmer, "modulate:a", 1.0, 0.25)
	tween.tween_property(card, "modulate:a", 1.0, 0.2)
	tween.tween_property(card, "scale", Vector2.ONE, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

func _make_box(color: Color, radius: int, shadow: int, padding: int) -> StyleBoxFlat:
	var box := StyleBoxFlat.new()
	box.bg_color = color
	box.set_corner_radius_all(radius)
	box.shadow_size = shadow
	box.shadow_offset = Vector2(0, 3.41)
	box.set_content_margin_all(padding)
	return box

func _return_to_main_lobby() -> void:
	if returning_to_lobby:
		return
	returning_to_lobby = true
	game_timer.stop()
	if game_socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		game_socket.close()
	get_tree().change_scene_to_file(MAIN_LOBBY_SCENE)
