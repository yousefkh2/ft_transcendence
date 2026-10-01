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
	if typeof(message) == TYPE_STRING:
		var response_msg = str_to_var(message)
		
		if(response_msg.op):
			print("Process message op: %s" % response_msg.op)
			
			if(response_msg.op == REQUEST_MATCH):
				print("REQUEST_MATCH")
		
		elif(response_msg.op == MATCH_PLAYERS):
			print("MATCH_PLAYERS")
			_enter_match_lobby(response_msg.response)
			
		elif(response_msg.op == PLAYER_JOINED):
			print("Player joined")
			var match_with_players = response_msg.response
			_build_player_lobby_list(match_with_players.users)
			
		elif(response_msg.op == PLAYER_LEFT):
			print("PLayer left")
			var match_with_players = response_msg.response
			print
			_build_player_lobby_list(match_with_players.users)
			
			
			

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
		if player is Dictionary:
			player_label.text = str(player.get("username", player.get("userId", "Unknown player")))
		else:
			player_label.text = str(player)
		player_label.custom_minimum_size = Vector2(260.0, 60.0)
		player_label.add_theme_font_override("font", preload("res://game_files/fonts/GrapeSoda.ttf"))
		player_label.add_theme_font_size_override("font_size", 48)
		player_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		
		if(player.role == NO_ROLE_ASSIGN):
			$Lobby_state/no_role_player.add_child(player_label)
		elif(player.role == ON_SITE_ASSIGN):
			$Lobby_state/on_site_player.add_child(player_label)
		elif(player.role == MISSION_CONTROL_ASSIGN):
			$Lobby_state/controll_layer.add_child(player_label)
		else:
			print("player has not been assinged to a role!")
		
	
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
	pass # Replace with function body.


func _on_join_on_site_button_pressed() -> void:
	pass # Replace with function body.


func _on_join_mission_control_button_pressed() -> void:
	pass # Replace with function body.
