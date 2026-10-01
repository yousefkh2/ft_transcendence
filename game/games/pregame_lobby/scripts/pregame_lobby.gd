extends Control

@onready var lobby_code: Label = $lobby_code_interface/lobby_code
@onready var game_type: Label = $create_lobby_interface/game_type
@onready var lobby_lang: Label = $create_lobby_interface/lobby_lang
@onready var user_name: Label = $player_data/User_name

var websocket_url = "ws://localhost:8080/ws/lobby"
var message_to_send = "TEST TEST TEST"

const REQUEST_MATCH = "REQUEST_MATCH"
const JOIN_MATCH = "JOIN_MATCH"
const MATCH_PLAYERS = "MATCH_PLAYERS"
const PLAYER_LEFT = "PLAYER_LEFT"
const PLAYER_JOINED = "PLAYER_JOINED"
const START_MATCH = "START_MATCH"
const MATCH_READY = "MATCH_READY"
const NO_ROLE_ASSIGN = "NO_ROLE_ASSIGN"
const MISSION_CONTROL_ASSIGN = "MISSION_CONTROL_ASSIGN"
const ON_SITE_ASSIGN = "ON_SITE_ASSIGN"
const CHECK_MATCH_READY = "CHECK_MATCH_READY"

const NO_ROLE_PLAYER = 0
const ON_SITE_PLAYER = 1
const MISSION_CONTROL_PLAYER = 2

@onready var _client : web_socket_client = $web_socket_client

func _ready() -> void:
	lobby_code.text = GameState.lobby_code
	game_type.text = GameState.game_mode
	print(GameState.player_name)
	print(GameState.lobby_data)
	user_name.text = GameState.player_name
	print(GameState.lobby_data)
	print(GameState.game_lang)
	lobby_lang.text = GameState.game_lang
	_build_player_lobby_list([GameState.player_name])
	print("Attemting to connect to server...")
	
	_connect_to_matchmaking_server()

	

func _send_message(message_to_send):
	var json_message = JSON.stringify(message_to_send)
	_client.send(json_message)

func _connect_to_matchmaking_server():
	var error = _client.connect_to_url(websocket_url)
	if (error != OK):
		print("Error connecting to websocket: %s " % [websocket_url])

func _process_received_message(message):
	if typeof(message) != TYPE_STRING:
		return

	var response_msg = JSON.parse_string(message)
	if typeof(response_msg) != TYPE_DICTIONARY:
		print("Invalid WebSocket message: %s" % message)
		return

	match response_msg.get("type", ""):
		"lobby.updated":
			_update_lobby_state(response_msg)
		"error":
			print("Server error: %s" % response_msg.get("message", "unknown error"))
		_:
			print("Unhandled server message: %s" % response_msg.get("type", ""))

func _enter_match_lobby(match_with_players):
	print("enter match lobby")
	print(match_with_players)
	
	_build_player_lobby_list(match_with_players.user)
	
	var match_id = match_with_players.matchInfo.matchId
	var check_match_ready = {
		"op": CHECK_MATCH_READY,
		"matchID": match_id
	}
	
	_send_message(check_match_ready)
	
	
func _update_lobby_state(response_msg: Dictionary) -> void:
	var players = response_msg.get("players", [])
	if players is Array:
		print("Lobby players updated: ", players)
		_build_player_lobby_list(players)
				

func _build_player_lobby_list(match_players):
	for team_child in $Lobby_state/on_site_player.get_children():
		team_child.queue_free()
		
	for team_child in $Lobby_state/controll_player.get_children():
		team_child.queue_free()
		
	for team_child in $Lobby_state/no_role_player.get_children():
		team_child.queue_free()
	
	for player in match_players:
		var player_label := Label.new()
		var role := ""
		if player is Dictionary:
			player_label.text = str(player.get("username", player.get("userId", "Unknown player")))
			role = str(player.get("role", "")).strip_edges().to_lower()
		else:
			player_label.text = str(player)
		player_label.custom_minimum_size = Vector2(260.0, 60.0)
		player_label.add_theme_font_override("font", preload("res://game_files/fonts/GrapeSoda.ttf"))
		player_label.add_theme_font_size_override("font_size", 48)
		player_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		
		match role:
			"on_site":
				$Lobby_state/on_site_player.add_child(player_label)
			"mission_control":
				$Lobby_state/controll_player.add_child(player_label)
			_:
				$Lobby_state/no_role_player.add_child(player_label)
		
	
func _on_websocket_message_recieved(message):
	print("Message received: %s " % message)
	_process_received_message(message)
	
func _on_websocket_client_connection_close():
	var ws = _client.get_socket()
	print("Client disconnected with code %s, reason: %s" % [ws.get_close_code(), ws.get_close_reason()])

func _on_websocket_client_connected_to_server():
	print("Client connected to server")

	_send_message({
		"type": "lobby.subscribe",
		"roomCode": GameState.lobby_code,
		"token": GameState.auth_token,
	})
	

func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
 

func _process(delta: float) -> void:
	pass


func _on_send_websocket_message_pressed() -> void:
	print("send out message from client to server")
	var message = {
		"type": "lobby.subscribe",
		"roomCode": GameState.lobby_code,
		"token": GameState.auth_token,
	}
	var error = _client.send_json(message)
	print(error)
	


func _on_join_no_role_button_pressed() -> void:
	_select_role("")

func _on_join_on_site_button_pressed() -> void:
	_select_role("on_site")


func _on_join_mission_control_button_pressed() -> void:
	_select_role("mission_control")


func _select_role(role: String) -> void:
	if _client.get_socket().get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Cannot select role: lobby WebSocket is not connected")
		return

	_send_message({
		"type": "lobby.role",
		"roomCode": GameState.lobby_code,
		"role": role,
		"token": GameState.auth_token,
	})
