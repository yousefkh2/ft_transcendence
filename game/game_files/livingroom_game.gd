extends Control

## Shared by both living room scenes (on site and mission control).
##
## Flow:
##  - Both players build the same puzzle from the lobby code (LivingroomPuzzle).
##  - On site sends its furniture layout ("game.layout") after joining and after
##    every swap; mission control asks for it on join ("game.layout_request").
##  - Mission control presses Check, compares the latest layout with the target
##    and sends the result ("game.check_result") to on site.
##  - The server only relays these messages between the two players.

const MAIN_LOBBY_SCENE := "res://games/lobby/lobby.tscn"
const STARTING_SECONDS := 5 * 60
const MAX_ATTEMPTS := 5
const GAME_WEBSOCKET_URL := "ws://localhost:8080/ws"
const GAME_FONT: FontFile = preload("res://game_files/fonts/GrapeSoda.ttf")
const ACCENT_RED := Color(0.7529412, 0.22745098, 0.2, 1)
const ACCENT_YELLOW := Color(1, 0.76862746, 0.20392157, 1)
const ACCENT_GREEN := Color(0.23, 0.6, 0.27, 1)
const INSTRUCTION_FONT_SIZE := 60
const INSTRUCTION_MIN_FONT_SIZE := 32
const TOAST_SECONDS := 2.5

@onready var game_timer: Timer = $GameTimer
@onready var timer_label: Label = $TimerLabel
@onready var leave_button: Button = $LeaveButton

var game_socket := WebSocketPeer.new()
var remaining_seconds := STARTING_SECONDS
var end_screen: Control
var returning_to_lobby := false
var game_room_join_sent := false

# Which side of the game this scene is; the scene's nodes tell us.
var is_on_site := false
var is_mission_control := false
# Mission control: the on-site layout as last received over the network.
var partner_layout: Array = []
var toast: PanelContainer
# Checks left this round. Mission control counts them and tells on site.
var attempts_left := MAX_ATTEMPTS
var attempts_label: Label
var toast_tween: Tween

func _ready() -> void:
	is_on_site = has_node("GridContainer")
	is_mission_control = has_node("first_instruction")
	_setup_puzzle()
	_create_attempts_label()
	game_timer.timeout.connect(_on_game_timer_timeout)
	_update_timer_label()
	game_timer.start()
	var error := game_socket.connect_to_url(GAME_WEBSOCKET_URL)
	if error != OK:
		push_error("Failed to connect to the game server: %s" % error)

# --- Puzzle setup -----------------------------------------------------------

func _setup_puzzle() -> void:
	if is_on_site:
		_place_start_furniture()
		for cell in $GridContainer.get_children():
			cell.swapped.connect(_send_layout)
	if is_mission_control:
		_show_instructions()
		$check_button.pressed.connect(_on_check_button_pressed)

func _place_start_furniture() -> void:
	var layout := LivingroomPuzzle.start_layout(GameState.lobby_code)
	var cells := $GridContainer.get_children()
	for i in cells.size():
		var item: String = layout[i]
		cells[i].icon = LivingroomPuzzle.ITEM_TEXTURES.get(item)

func _show_instructions() -> void:
	# Fill the labels top to bottom, whatever their node names are.
	var labels: Array[Label] = [$first_instruction, $second_instruction, $third_instruction, $fourth_instruction]
	labels.sort_custom(func(a: Label, b: Label): return a.position.y < b.position.y)
	var clues := LivingroomPuzzle.clues(GameState.lobby_code)
	for i in labels.size():
		labels[i].text = LivingroomPuzzle.clue_text(clues[i]) if i < clues.size() else ""
		_fit_label_text(labels[i])

# Shrink long clues (some languages are much longer) so they stay on one line.
func _fit_label_text(label: Label) -> void:
	var font := label.get_theme_font("font")
	var text := label.text.to_upper() if label.uppercase else label.text
	var font_size := INSTRUCTION_FONT_SIZE
	while font_size > INSTRUCTION_MIN_FONT_SIZE and font.get_string_size(text, HORIZONTAL_ALIGNMENT_LEFT, -1, font_size).x > label.size.x:
		font_size -= 2
	label.add_theme_font_size_override("font_size", font_size)

## On site: the current furniture layout, read from the grid's icons.
func _current_layout() -> Array:
	var layout := []
	for cell in $GridContainer.get_children():
		layout.append(LivingroomPuzzle.ITEM_TEXTURES.find_key(cell.icon) if cell.icon else "")
	return layout

# --- Checking ----------------------------------------------------------------

func _on_check_button_pressed() -> void:
	if end_screen or attempts_left <= 0:
		return
	if partner_layout.is_empty():
		_send_game_message("game.layout_request")
		_show_toast(tr("_GAME_CHECK_WAITING_"))
		return

	var target := LivingroomPuzzle.target_layout(GameState.lobby_code)
	var correct := 0
	for item in LivingroomPuzzle.ITEMS:
		if partner_layout.find(item) == target.find(item):
			correct += 1
	var solved := correct == LivingroomPuzzle.ITEMS.size()
	if not solved:
		attempts_left -= 1

	_send_game_message("game.check_result", {"solved": solved, "correct": correct, "attemptsLeft": attempts_left})
	_show_check_result(solved, correct, attempts_left)

## Shown on both screens: mission control right after checking, on site when
## the result arrives over the network.
func _show_check_result(solved: bool, correct: int, attempts: int) -> void:
	attempts_left = attempts
	_update_attempts_label()
	if solved:
		_show_end_screen(tr("_GAME_ROOM_DONE_"), tr("_GAME_ROOM_DONE_TEXT_"), ACCENT_GREEN)
	elif attempts_left <= 0:
		_show_end_screen(tr("_GAME_NO_ATTEMPTS_"), tr("_GAME_NO_ATTEMPTS_TEXT_").format({"max": MAX_ATTEMPTS}), ACCENT_RED)
	else:
		var details := tr("_GAME_CHECK_WRONG_").format({"n": correct, "total": LivingroomPuzzle.ITEMS.size()})
		details += "\n" + tr("_GAME_ATTEMPTS_LEFT_").format({"left": attempts_left})
		_show_toast(details, tr("_GAME_CHECK_WRONG_TITLE_"))

# --- Networking --------------------------------------------------------------

func _process(_delta: float) -> void:
	game_socket.poll()
	if game_socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		return
	if not game_room_join_sent:
		_join_game_room()
		game_room_join_sent = true

	while game_socket.get_available_packet_count() > 0:
		var response = JSON.parse_string(game_socket.get_packet().get_string_from_utf8())
		if typeof(response) == TYPE_DICTIONARY:
			_handle_server_message(response)

func _handle_server_message(response: Dictionary) -> void:
	var payload = response.get("payload", {})
	if typeof(payload) != TYPE_DICTIONARY:
		payload = {}
	match response.get("type", ""):
		"room.joined":
			if is_on_site:
				_send_layout()
			else:
				_send_game_message("game.layout_request")
		"game.layout_request":
			if is_on_site:
				_send_layout()
		"game.layout":
			var layout = payload.get("layout", [])
			if is_mission_control and typeof(layout) == TYPE_ARRAY and layout.size() == LivingroomPuzzle.SPOT_COUNT:
				partner_layout = layout
		"game.check_result":
			if is_on_site and not end_screen:
				_show_check_result(bool(payload.get("solved", false)), int(payload.get("correct", 0)), int(payload.get("attemptsLeft", attempts_left)))
		"player.disconnected":
			# Once the round is over (time's up or solved) let each player
			# leave the end screen on their own instead of pulling both out.
			if end_screen:
				return
			GameState.lobby_message = str(response.get("message", "A player left the game"))
			_return_to_main_lobby()

func _join_game_room() -> void:
	game_socket.send_text(JSON.stringify({
		"type": "room.join",
		"roomCode": GameState.lobby_code,
		"token": GameState.auth_token,
	}))

func _send_layout() -> void:
	_send_game_message("game.layout", {"layout": _current_layout()})

func _send_game_message(type: String, payload: Dictionary = {}) -> void:
	if game_socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		return
	game_socket.send_text(JSON.stringify({"type": type, "payload": payload}))

# --- Timer -------------------------------------------------------------------

func _on_game_timer_timeout() -> void:
	remaining_seconds -= 1
	_update_timer_label()
	if remaining_seconds <= 0:
		_show_end_screen(tr("_GAME_TIMES_UP_"), tr("_GAME_TIMES_UP_TEXT_"), ACCENT_RED)

func _update_timer_label() -> void:
	timer_label.text = "%d:%02d" % [remaining_seconds / 60, remaining_seconds % 60]

func _on_leave_button_pressed() -> void:
	_return_to_main_lobby()

# --- UI ----------------------------------------------------------------------

## Short message card at the top of the screen that fades out by itself,
## with an optional red title line above the text.
func _show_toast(text: String, title_text: String = "") -> void:
	if not toast:
		toast = PanelContainer.new()
		toast.mouse_filter = Control.MOUSE_FILTER_IGNORE
		toast.add_theme_stylebox_override("panel", _make_box(Color(1, 1, 1, 1), 16, 5, 20))
		var lines := VBoxContainer.new()
		lines.add_theme_constant_override("separation", 6)
		toast.add_child(lines)
		var title := _make_label(56, ACCENT_RED)
		title.add_theme_color_override("font_shadow_color", Color(0.07058824, 0.07058824, 0.07058824, 1))
		title.add_theme_constant_override("shadow_offset_x", 3)
		title.add_theme_constant_override("shadow_offset_y", 3)
		lines.add_child(title)
		lines.add_child(_make_label(40, Color(0, 0, 0, 1)))
		add_child(toast)
	var toast_title: Label = toast.get_child(0).get_child(0)
	toast_title.text = title_text
	toast_title.visible = title_text != ""
	toast.get_child(0).get_child(1).text = text
	toast.reset_size()
	toast.position = Vector2((size.x - toast.size.x) / 2.0, 40)
	toast.modulate.a = 1.0
	if toast_tween:
		toast_tween.kill()
	toast_tween = create_tween()
	toast_tween.tween_interval(TOAST_SECONDS)
	toast_tween.tween_property(toast, "modulate:a", 0.0, 0.4)

func _make_label(font_size: int, color: Color) -> Label:
	var label := Label.new()
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.add_theme_font_override("font", GAME_FONT)
	label.add_theme_font_size_override("font_size", font_size)
	label.add_theme_color_override("font_color", color)
	return label

## "Attempts left" counter under the timer, on both screens.
func _create_attempts_label() -> void:
	attempts_label = _make_label(36, Color(1, 1, 1, 1))
	attempts_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_LEFT
	attempts_label.add_theme_color_override("font_shadow_color", Color(0.07058824, 0.07058824, 0.07058824, 1))
	attempts_label.add_theme_constant_override("shadow_offset_x", 3)
	attempts_label.add_theme_constant_override("shadow_offset_y", 3)
	attempts_label.position = timer_label.position + Vector2(0, timer_label.size.y)
	add_child(attempts_label)
	_update_attempts_label()

func _update_attempts_label() -> void:
	attempts_label.text = tr("_GAME_ATTEMPTS_LEFT_").format({"left": attempts_left})

## Full-screen end card (time's up / room complete) with a back-to-lobby button.
func _show_end_screen(title_text: String, subtitle_text: String, title_color: Color) -> void:
	if end_screen:
		return
	game_timer.stop()
	leave_button.disabled = true
	if toast:
		toast.hide()

	# Full-screen dimmer that also blocks clicks/drags on the game underneath.
	end_screen = Control.new()
	end_screen.set_anchors_preset(Control.PRESET_FULL_RECT)
	end_screen.mouse_filter = Control.MOUSE_FILTER_STOP

	var dimmer := ColorRect.new()
	dimmer.color = Color(0, 0, 0, 0.65)
	dimmer.set_anchors_preset(Control.PRESET_FULL_RECT)
	dimmer.mouse_filter = Control.MOUSE_FILTER_IGNORE
	end_screen.add_child(dimmer)

	var center := CenterContainer.new()
	center.set_anchors_preset(Control.PRESET_FULL_RECT)
	center.mouse_filter = Control.MOUSE_FILTER_IGNORE
	end_screen.add_child(center)

	var card := PanelContainer.new()
	card.custom_minimum_size = Vector2(760, 0)
	card.add_theme_stylebox_override("panel", _make_box(Color(1, 1, 1, 1), 24, 12, 48))
	center.add_child(card)

	var contents := VBoxContainer.new()
	contents.alignment = BoxContainer.ALIGNMENT_CENTER
	contents.add_theme_constant_override("separation", 28)
	card.add_child(contents)

	var title := Label.new()
	title.text = title_text
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_override("font", GAME_FONT)
	title.add_theme_font_size_override("font_size", 110)
	title.add_theme_color_override("font_color", title_color)
	title.add_theme_color_override("font_shadow_color", Color(0.07058824, 0.07058824, 0.07058824, 1))
	title.add_theme_constant_override("shadow_offset_x", 4)
	title.add_theme_constant_override("shadow_offset_y", 4)
	contents.add_child(title)

	var subtitle := Label.new()
	subtitle.text = subtitle_text
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

	add_child(end_screen)

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
